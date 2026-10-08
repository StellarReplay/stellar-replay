// Package capture records one supported Stellar JSON-RPC interaction.
package capture

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
)

const (
	DefaultTimeout          = 30 * time.Second
	DefaultMaxResponseBytes = fixture.MaxFixtureBytes
)

var (
	ErrUnsupportedMethod = errors.New("unsupported RPC method")
	ErrRedirect          = errors.New("redirects are disabled")
)

// Config defines one explicit, read-only capture operation.
type Config struct {
	Endpoint         string
	Network          string
	ToolVersion      string
	Method           string
	ID               json.RawMessage
	Params           json.RawMessage
	Timeout          time.Duration
	MaxResponseBytes int64

	// Transport is intended for tests and controlled callers. When nil, capture
	// uses a transport with no proxy and private-destination protection.
	Transport http.RoundTripper
	// Now is injectable so tests can prove deterministic fixture output.
	Now func() time.Time
	// Sanitize may redact or reject response content before it is sealed.
	Sanitize func(*fixture.Response) error
}

// Capture performs exactly one HTTP POST and returns a sealed fixture. It never
// retries, follows redirects, or falls back to another endpoint.
func Capture(ctx context.Context, cfg Config) (fixture.Fixture, error) {
	var captured fixture.Fixture
	if ctx == nil {
		return captured, errors.New("capture context is nil")
	}
	if !fixture.Supported(cfg.Method) {
		return captured, fmt.Errorf("%w %q", ErrUnsupportedMethod, cfg.Method)
	}
	endpoint, err := validateEndpoint(cfg.Endpoint)
	if err != nil {
		return captured, err
	}
	if cfg.Network == "" || cfg.ToolVersion == "" {
		return captured, errors.New("network and toolVersion are required")
	}
	if len(cfg.ID) == 0 {
		return captured, errors.New("request id is required")
	}
	if !scalarJSON(cfg.ID) {
		return captured, errors.New("request id must be a scalar JSON value")
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	if timeout < 0 {
		return captured, errors.New("timeout must be positive")
	}
	maxBytes := cfg.MaxResponseBytes
	if maxBytes == 0 {
		maxBytes = DefaultMaxResponseBytes
	}
	if maxBytes <= 0 || maxBytes > fixture.MaxFixtureBytes {
		return captured, fmt.Errorf("max response bytes must be between 1 and %d", fixture.MaxFixtureBytes)
	}
	now := cfg.Now
	if now == nil {
		now = time.Now
	}

	payload := requestPayload{JSONRPC: "2.0", ID: cfg.ID, Method: cfg.Method}
	if len(cfg.Params) > 0 {
		if !json.Valid(cfg.Params) {
			return captured, errors.New("request params must be valid JSON")
		}
		payload.Params = cfg.Params
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return captured, fmt.Errorf("marshal request: %w", err)
	}

	requestCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return captured, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{
		Transport: cfg.Transport,
		Timeout:   timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return ErrRedirect
		},
	}
	if client.Transport == nil {
		client.Transport = secureTransport()
	}
	response, err := client.Do(req)
	if err != nil {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
		if errors.Is(err, ErrRedirect) {
			return captured, ErrRedirect
		}
		return captured, fmt.Errorf("capture request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusMultipleChoices && response.StatusCode < http.StatusBadRequest {
		return captured, fmt.Errorf("%w: HTTP status %d", ErrRedirect, response.StatusCode)
	}
	responseBody, err := readBounded(response.Body, maxBytes)
	if err != nil {
		return captured, err
	}
	var rpcResponse fixture.Response
	decoder := json.NewDecoder(bytes.NewReader(responseBody))
	if err := decoder.Decode(&rpcResponse); err != nil {
		return captured, fmt.Errorf("decode RPC response: %w", err)
	}
	if err := ensureEOF(decoder); err != nil {
		return captured, err
	}
	if cfg.Sanitize != nil {
		if err := cfg.Sanitize(&rpcResponse); err != nil {
			return captured, fmt.Errorf("sanitize response: %w", err)
		}
	}

	captured = fixture.Fixture{
		SchemaVersion: fixture.SchemaVersion,
		CapturedAt:    now().UTC().Format(time.RFC3339),
		Provenance: fixture.Provenance{
			Endpoint: endpoint.String(), Network: cfg.Network, ToolVersion: cfg.ToolVersion,
		},
		Request:  fixture.Request{JSONRPC: payload.JSONRPC, ID: cfg.ID, Method: cfg.Method, Params: cfg.Params},
		Response: rpcResponse,
	}
	if err := fixture.Seal(&captured); err != nil {
		return fixture.Fixture{}, fmt.Errorf("validate captured fixture: %w", err)
	}
	return captured, nil
}

// Record captures one interaction and atomically writes the sealed fixture.
func Record(ctx context.Context, cfg Config, path string) (fixture.Fixture, error) {
	captured, err := Capture(ctx, cfg)
	if err != nil {
		return fixture.Fixture{}, err
	}
	if err := fixture.Save(path, captured); err != nil {
		return fixture.Fixture{}, fmt.Errorf("save captured fixture: %w", err)
	}
	return captured, nil
}

type requestPayload struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func validateEndpoint(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, errors.New("endpoint is required")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return nil, errors.New("endpoint must be an HTTPS URL without userinfo")
	}
	if u.Fragment != "" {
		return nil, errors.New("endpoint must not contain a URL fragment")
	}
	if strings.ContainsAny(u.Host, "\r\n") {
		return nil, errors.New("endpoint contains invalid host characters")
	}
	host := u.Hostname()
	if ip := net.ParseIP(host); ip != nil && disallowedIP(ip) {
		return nil, errors.New("endpoint must not target a private, loopback, or link-local address")
	}
	for key := range u.Query() {
		if strings.Contains(strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(key)), "token") || strings.Contains(strings.ToLower(key), "secret") {
			return nil, fmt.Errorf("endpoint contains sensitive query field %q", key)
		}
	}
	return u, nil
}

func secureTransport() http.RoundTripper {
	return &http.Transport{
		Proxy:           nil,
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		DialContext:     secureDialContext,
	}
}

func secureDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	if ip := net.ParseIP(host); ip != nil {
		if disallowedIP(ip) {
			return nil, errors.New("destination is private, loopback, or link-local")
		}
		return (&net.Dialer{}).DialContext(ctx, network, net.JoinHostPort(host, port))
	}
	addresses, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("resolve endpoint: %w", err)
	}
	dialer := &net.Dialer{}
	var lastErr error
	for _, address := range addresses {
		if disallowedIP(address.IP) {
			continue
		}
		conn, dialErr := dialer.DialContext(ctx, network, net.JoinHostPort(address.IP.String(), port))
		if dialErr == nil {
			return conn, nil
		}
		lastErr = dialErr
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("endpoint resolved only to private, loopback, or link-local addresses")
}

func disallowedIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified()
}

func readBounded(reader io.Reader, max int64) ([]byte, error) {
	if max > math.MaxInt32 {
		max = math.MaxInt32
	}
	data, err := io.ReadAll(io.LimitReader(reader, max+1))
	if err != nil {
		return nil, fmt.Errorf("read RPC response: %w", err)
	}
	if int64(len(data)) > max {
		return nil, fmt.Errorf("RPC response exceeds %d bytes", max)
	}
	return data, nil
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return fmt.Errorf("trailing JSON: %w", err)
	}
	return nil
}

func scalarJSON(raw json.RawMessage) bool {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if decoder.Decode(&value) != nil || ensureEOF(decoder) != nil {
		return false
	}
	switch value.(type) {
	case nil, string, json.Number:
		return true
	default:
		return false
	}
}

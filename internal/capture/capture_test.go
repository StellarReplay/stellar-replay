package capture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
)

func TestCaptureWithLocalTLSServerAndRecord(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected request: %s %s %s", r.Method, r.URL, r.Header.Get("Content-Type"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), `"method":"getHealth"`) {
			t.Errorf("request body did not contain method: %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":{"status":"healthy"}}`)
	}))
	defer server.Close()

	path := t.TempDir() + "\\health.json"
	cfg := testConfig(server, "getHealth", json.RawMessage(`{}`))
	first, err := Record(context.Background(), cfg, path)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := fixture.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if first.Integrity.SHA256 != loaded.Integrity.SHA256 || calls.Load() != 1 {
		t.Fatalf("record mismatch: hash=%v calls=%d", loaded.Integrity, calls.Load())
	}
}

func TestCaptureIsDeterministicWithFixedClock(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":7,"result":{"networkPassphrase":"test"}}`)
	}))
	defer server.Close()
	cfg := testConfig(server, "getNetwork", nil)
	cfg.ID = json.RawMessage(`7`)
	cfg.Now = func() time.Time { return time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC) }
	first, err := Capture(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Capture(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if first.Integrity.SHA256 != second.Integrity.SHA256 {
		t.Fatalf("fixed-clock captures differ: %s != %s", first.Integrity.SHA256, second.Integrity.SHA256)
	}
}

func TestCapturePreservesRPCErrorAndSanitizesResponse(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(r, 200, `{"jsonrpc":"2.0","id":2,"error":{"code":-32602,"message":"invalid params","data":{"secret":"remove"}}}`), nil
	})
	cfg := testConfigWithTransport(transport, "getLedgerEntries", json.RawMessage(`{"keys":["AQIDBA=="]}`))
	cfg.ID = json.RawMessage(`2`)
	cfg.Sanitize = func(response *fixture.Response) error {
		response.Error.Data = json.RawMessage(`{"redacted":true}`)
		return nil
	}
	captured, err := Capture(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if captured.Response.Error == nil || captured.Response.Error.Code != -32602 || string(captured.Response.Error.Data) != `{"redacted":true}` {
		t.Fatalf("RPC error was not preserved/sanitized: %+v", captured.Response.Error)
	}
}

func TestCaptureRejectsUnsupportedMethodWithoutContactingEndpoint(t *testing.T) {
	var calls atomic.Int32
	transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, errors.New("unexpected network call")
	})
	cfg := testConfigWithTransport(transport, "sendTransaction", nil)
	_, err := Capture(context.Background(), cfg)
	if err == nil || !errors.Is(err, ErrUnsupportedMethod) || calls.Load() != 0 {
		t.Fatalf("unexpected result: err=%v calls=%d", err, calls.Load())
	}
}

func TestCaptureRejectsMalformedOversizedAndRedirectResponses(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		want     string
	}{
		{"malformed", jsonResponse(nil, 200, `{not-json`), "decode RPC response"},
		{"oversized", jsonResponse(nil, 200, strings.Repeat("x", 32)), "exceeds 16 bytes"},
		{"redirect", redirectResponse(), "redirects are disabled"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
				response := test.response
				response.Request = r
				return response, nil
			})
			cfg := testConfigWithTransport(transport, "getHealth", json.RawMessage(`{}`))
			cfg.MaxResponseBytes = 16
			_, err := Capture(context.Background(), cfg)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected %q, got %v", test.want, err)
			}
		})
	}
}

func TestCaptureTimeoutAndNoRetry(t *testing.T) {
	var calls atomic.Int32
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		<-r.Context().Done()
		return nil, r.Context().Err()
	})
	cfg := testConfigWithTransport(transport, "getHealth", json.RawMessage(`{}`))
	cfg.Timeout = 20 * time.Millisecond
	_, err := Capture(context.Background(), cfg)
	if err == nil || calls.Load() != 1 || !strings.Contains(err.Error(), "capture request") {
		t.Fatalf("expected one timed-out request, err=%v calls=%d", err, calls.Load())
	}
}

func TestCaptureEndpointPolicy(t *testing.T) {
	for _, endpoint := range []string{
		"http://rpc.example.test",
		"https://user:pass@rpc.example.test",
		"https://127.0.0.1/rpc",
		"https://rpc.example.test/rpc?api_key=secret",
		"https://rpc.example.test/rpc?token=secret-value",
		"https://rpc.example.test/rpc?X-Token=secret-value",
		"https://rpc.example.test/rpc?oauth_token=secret-value",
		"https://rpc.example.test/rpc?ID_TOKEN=secret-value",
	} {
		cfg := testConfigWithTransport(roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("unexpected network call")
		}), "getHealth", json.RawMessage(`{}`))
		cfg.Endpoint = endpoint
		if _, err := Capture(context.Background(), cfg); err == nil {
			t.Fatalf("endpoint %q was accepted", endpoint)
		} else if strings.Contains(err.Error(), "secret-value") {
			t.Fatalf("endpoint error exposed a credential value: %v", err)
		}
	}
}

func TestCaptureRejectsMismatchedResponseID(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(r, http.StatusOK, `{"jsonrpc":"2.0","id":2,"result":{"status":"healthy"}}`), nil
	})
	_, err := Capture(context.Background(), testConfigWithTransport(transport, "getHealth", nil))
	if err == nil || !strings.Contains(err.Error(), "IDs must match") {
		t.Fatalf("mismatched response ID was accepted: %v", err)
	}
}

func TestCaptureRejectsOversizedRequestBeforeTransport(t *testing.T) {
	var calls atomic.Int32
	cfg := testConfigWithTransport(roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, errors.New("unexpected network call")
	}), "getHealth", json.RawMessage(`{"padding":"`+strings.Repeat("x", 64)+`"}`))
	cfg.MaxRequestBytes = 32
	if _, err := Capture(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "request exceeds") || calls.Load() != 0 {
		t.Fatalf("oversized request was not rejected before transport: err=%v calls=%d", err, calls.Load())
	}
}

func TestCaptureDefaultSanitizationRedactsResponseSecrets(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(r, http.StatusOK, `{"jsonrpc":"2.0","id":1,"result":{"api_key":"secret","value":"safe"}}`), nil
	})
	cfg := testConfigWithTransport(transport, "getHealth", nil)
	captured, err := Capture(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(captured.Response.Result), "secret") || !strings.Contains(string(captured.Response.Result), "[REDACTED]") {
		t.Fatalf("response secret was stored: %s", captured.Response.Result)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return fn(r) }

func testConfig(server *httptest.Server, method string, params json.RawMessage) Config {
	parsed, err := url.Parse(server.URL)
	if err != nil {
		panic(err)
	}
	return Config{
		Endpoint: "https://rpc.example.test/rpc",
		Network:  "testnet", ToolVersion: "test", Method: method, ID: json.RawMessage(`1`), Params: params,
		Transport: rewriteTransport{server: server, target: parsed},
		Now:       func() time.Time { return time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC) },
	}
}

func testConfigWithTransport(transport http.RoundTripper, method string, params json.RawMessage) Config {
	return Config{
		Endpoint: "https://rpc.example.test/rpc",
		Network:  "testnet", ToolVersion: "test", Method: method, ID: json.RawMessage(`1`), Params: params,
		Transport: transport,
		Now:       func() time.Time { return time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC) },
	}
}

type rewriteTransport struct {
	server *httptest.Server
	target *url.URL
}

func (transport rewriteTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	urlCopy := *transport.target
	urlCopy.Path = request.URL.Path
	urlCopy.RawQuery = request.URL.RawQuery
	clone.URL = &urlCopy
	clone.Host = urlCopy.Host
	return transport.server.Client().Transport.RoundTrip(clone)
}

func jsonResponse(request *http.Request, status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status, Status: fmt.Sprintf("%d", status), Body: io.NopCloser(strings.NewReader(body)),
		Header: make(http.Header), Request: request,
	}
}

func redirectResponse() *http.Response {
	response := jsonResponse(nil, http.StatusFound, "")
	response.Header.Set("Location", "https://rpc.example.test/other")
	return response
}

// Package fixture defines the versioned Stellar Replay fixture contract.
package fixture

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	SchemaVersion   = "1"
	MaxFixtureBytes = 8 << 20
	MaxLedgerKeys   = 200
)

var supportedMethods = map[string]struct{}{
	"getHealth": {}, "getLatestLedger": {}, "getNetwork": {}, "getLedgerEntries": {},
}

// Fixture is one captured JSON-RPC request/response interaction.
type Fixture struct {
	SchemaVersion string     `json:"schemaVersion"`
	CapturedAt    string     `json:"capturedAt"`
	Provenance    Provenance `json:"provenance"`
	Request       Request    `json:"request"`
	Response      Response   `json:"response"`
	Integrity     *Integrity `json:"integrity,omitempty"`
}

type Provenance struct {
	Endpoint     string          `json:"endpoint"`
	Network      string          `json:"network"`
	ToolVersion  string          `json:"toolVersion"`
	ChainContext json.RawMessage `json:"chainContext,omitempty"`
}

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Integrity struct {
	Algorithm string `json:"algorithm"`
	SHA256    string `json:"sha256"`
}

// Supported reports whether method is frozen into the v0.1 fixture contract.
func Supported(method string) bool {
	_, ok := supportedMethods[method]
	return ok
}

// ValidateRequest checks the supported JSON-RPC request contract without
// requiring a response or fixture provenance.
func ValidateRequest(r Request) error {
	return validateRequest(r)
}

// CanonicalJSONValue returns deterministic JSON for one arbitrary JSON value.
// It is used by other internal packages that share fixture matching semantics.
func CanonicalJSONValue(raw json.RawMessage) ([]byte, error) {
	return canonicalJSON(raw)
}

// Validate checks the schema, supported request contract, and integrity value.
func Validate(f Fixture) error {
	if err := validateStructure(f); err != nil {
		return err
	}
	if f.Integrity == nil {
		return errors.New("integrity is required")
	}
	if f.Integrity.Algorithm != "sha256" || !isSHA256(f.Integrity.SHA256) {
		return errors.New("integrity must use a lowercase hexadecimal sha256 value")
	}
	return VerifyIntegrity(f)
}

// Seal computes and attaches the integrity hash after validating the fixture.
func Seal(f *Fixture) error {
	if f == nil {
		return errors.New("fixture is nil")
	}
	f.Integrity = nil
	if err := validateStructure(*f); err != nil {
		return err
	}
	hash, err := Hash(*f)
	if err != nil {
		return err
	}
	f.Integrity = &Integrity{Algorithm: "sha256", SHA256: hash}
	return nil
}

// Hash returns SHA-256 over canonical fixture JSON without the integrity field.
func Hash(f Fixture) (string, error) {
	if err := validateStructure(f); err != nil {
		return "", err
	}
	f.Integrity = nil
	encoded, err := canonicalFixture(f)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

// VerifyIntegrity verifies the stored hash without changing the fixture.
func VerifyIntegrity(f Fixture) error {
	if f.Integrity == nil {
		return errors.New("integrity is required")
	}
	hash, err := Hash(f)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare([]byte(hash), []byte(f.Integrity.SHA256)) != 1 {
		return errors.New("fixture integrity hash mismatch")
	}
	return nil
}

// CanonicalJSON returns deterministic JSON, including integrity when present.
func CanonicalJSON(f Fixture) ([]byte, error) {
	if err := validateStructure(f); err != nil {
		return nil, err
	}
	return canonicalFixture(f)
}

// Save seals and atomically writes a fixture with restrictive permissions.
func Save(path string, f Fixture) error {
	if path == "" {
		return errors.New("fixture path is required")
	}
	if err := Seal(&f); err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal fixture: %w", err)
	}
	encoded = append(encoded, '\n')
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create fixture directory: %w", err)
	}
	temp, err := os.CreateTemp(dir, ".stellar-replay-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary fixture: %w", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if err := temp.Chmod(0o600); err != nil {
		_ = temp.Close()
		return fmt.Errorf("set fixture permissions: %w", err)
	}
	if _, err := temp.Write(encoded); err != nil {
		_ = temp.Close()
		return fmt.Errorf("write fixture: %w", err)
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return fmt.Errorf("sync fixture: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close fixture: %w", err)
	}
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("install fixture: %w", err)
	}
	return nil
}

// Load reads, strictly decodes, validates, and verifies a fixture.
func Load(path string) (Fixture, error) {
	var f Fixture
	if path == "" {
		return f, errors.New("fixture path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return f, fmt.Errorf("open fixture: %w", err)
	}
	defer file.Close()
	limited := io.LimitReader(file, MaxFixtureBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return f, fmt.Errorf("read fixture: %w", err)
	}
	if len(data) > MaxFixtureBytes {
		return f, fmt.Errorf("fixture exceeds %d bytes", MaxFixtureBytes)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&f); err != nil {
		return f, fmt.Errorf("decode fixture: %w", err)
	}
	if err := ensureEOF(decoder); err != nil {
		return f, err
	}
	if err := Validate(f); err != nil {
		return f, fmt.Errorf("validate fixture: %w", err)
	}
	return f, nil
}

func validateStructure(f Fixture) error {
	if f.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema version %q", f.SchemaVersion)
	}
	if _, err := time.Parse(time.RFC3339, f.CapturedAt); err != nil {
		return fmt.Errorf("capturedAt must be RFC3339: %w", err)
	}
	if err := validateProvenance(f.Provenance); err != nil {
		return err
	}
	if err := validateRequest(f.Request); err != nil {
		return err
	}
	if err := validateResponse(f.Response); err != nil {
		return err
	}
	if !sameScalarJSON(f.Request.ID, f.Response.ID) {
		return errors.New("request and response IDs must match")
	}
	return nil
}

func validateProvenance(p Provenance) error {
	if p.Endpoint == "" || p.Network == "" || p.ToolVersion == "" {
		return errors.New("provenance endpoint, network, and toolVersion are required")
	}
	u, err := url.Parse(p.Endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return errors.New("provenance endpoint must be an HTTPS URL without userinfo")
	}
	if err := rejectSensitiveURL(u); err != nil {
		return err
	}
	if len(p.ChainContext) > 0 {
		if _, err := canonicalJSON(p.ChainContext); err != nil {
			return fmt.Errorf("invalid chainContext: %w", err)
		}
	}
	return nil
}

func validateRequest(r Request) error {
	if r.JSONRPC != "2.0" || !isScalarJSON(r.ID) || r.Method == "" {
		return errors.New("request requires jsonrpc 2.0, scalar id, and method")
	}
	if !Supported(r.Method) {
		return fmt.Errorf("unsupported RPC method %q", r.Method)
	}
	if err := rejectSensitiveJSON(r.Params, "request params"); err != nil {
		return err
	}
	if r.Method == "getLedgerEntries" {
		return validateLedgerEntriesParams(r.Params)
	}
	if len(r.Params) == 0 {
		return nil
	}
	var params map[string]json.RawMessage
	if err := json.Unmarshal(r.Params, &params); err != nil || params == nil || len(params) != 0 {
		return fmt.Errorf("%s accepts only omitted params or an empty object", r.Method)
	}
	return nil
}

func validateLedgerEntriesParams(raw json.RawMessage) error {
	if len(raw) == 0 {
		return errors.New("getLedgerEntries params are required")
	}
	var params map[string]json.RawMessage
	if err := json.Unmarshal(raw, &params); err != nil || params == nil {
		return errors.New("getLedgerEntries params must be an object")
	}
	for key := range params {
		if key != "keys" && key != "xdrFormat" {
			return fmt.Errorf("getLedgerEntries has unsupported parameter %q", key)
		}
	}
	keysRaw, ok := params["keys"]
	if !ok {
		return errors.New("getLedgerEntries keys are required")
	}
	var keys []string
	if err := json.Unmarshal(keysRaw, &keys); err != nil || len(keys) == 0 {
		return errors.New("getLedgerEntries keys must be a non-empty array")
	}
	if len(keys) > MaxLedgerKeys {
		return fmt.Errorf("getLedgerEntries accepts at most %d keys", MaxLedgerKeys)
	}
	for _, key := range keys {
		decoded, err := base64.StdEncoding.DecodeString(key)
		if err != nil || len(decoded) == 0 {
			return errors.New("getLedgerEntries keys must be non-empty base64 strings")
		}
	}
	if format, ok := params["xdrFormat"]; ok {
		var value string
		if err := json.Unmarshal(format, &value); err != nil || (value != "base64" && value != "json") {
			return errors.New("xdrFormat must be base64 or json")
		}
	}
	return nil
}

func validateResponse(r Response) error {
	if r.JSONRPC != "2.0" || !isScalarJSON(r.ID) {
		return errors.New("response requires jsonrpc 2.0 and scalar id")
	}
	if (len(r.Result) == 0) == (r.Error == nil) {
		return errors.New("response must contain exactly one result or error")
	}
	if len(r.Result) > 0 {
		if _, err := canonicalJSON(r.Result); err != nil {
			return fmt.Errorf("invalid response result: %w", err)
		}
	}
	if r.Error != nil {
		if r.Error.Message == "" {
			return errors.New("response error message is required")
		}
		if len(r.Error.Data) > 0 {
			if _, err := canonicalJSON(r.Error.Data); err != nil {
				return fmt.Errorf("invalid response error data: %w", err)
			}
		}
	}
	return nil
}

func canonicalFixture(f Fixture) ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("marshal fixture: %w", err)
	}
	return canonicalJSON(data)
}

func canonicalJSON(data []byte) ([]byte, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("JSON value is required")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if err := ensureEOF(decoder); err != nil {
		return nil, err
	}
	return marshalCanonical(value)
}

func marshalCanonical(value any) ([]byte, error) {
	switch typed := value.(type) {
	case nil, bool, string, json.Number:
		return json.Marshal(typed)
	case []any:
		parts := make([][]byte, len(typed))
		for i, item := range typed {
			part, err := marshalCanonical(item)
			if err != nil {
				return nil, err
			}
			parts[i] = part
		}
		return joinJSON('[', ']', parts), nil
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		var out bytes.Buffer
		out.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				out.WriteByte(',')
			}
			encodedKey, _ := json.Marshal(key)
			valueJSON, err := marshalCanonical(typed[key])
			if err != nil {
				return nil, err
			}
			out.Write(encodedKey)
			out.WriteByte(':')
			out.Write(valueJSON)
		}
		out.WriteByte('}')
		return out.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported JSON value %T", value)
	}
}

func joinJSON(open, close byte, parts [][]byte) []byte {
	var out bytes.Buffer
	out.WriteByte(open)
	for i, part := range parts {
		if i > 0 {
			out.WriteByte(',')
		}
		out.Write(part)
	}
	out.WriteByte(close)
	return out.Bytes()
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

func isScalarJSON(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var value any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
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

func rejectSensitiveJSON(raw json.RawMessage, location string) error {
	if len(raw) == 0 {
		return nil
	}
	var value any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return fmt.Errorf("%s is invalid JSON: %w", location, err)
	}
	return scanSensitive(value, location)
}

func scanSensitive(value any, location string) error {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			if sensitiveKey(key) {
				return fmt.Errorf("%s contains prohibited sensitive field %q", location, key)
			}
			if err := scanSensitive(child, location+"."+key); err != nil {
				return err
			}
		}
	case []any:
		for i, child := range typed {
			if err := scanSensitive(child, fmt.Sprintf("%s[%d]", location, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func sensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("_", "", "-", "", " ", "").Replace(key))
	for _, prohibited := range []string{"privatekey", "secret", "seedphrase", "mnemonic", "password", "authorization", "token", "accesstoken", "refreshtoken", "apikey", "bearer"} {
		if strings.Contains(normalized, prohibited) {
			return true
		}
	}
	return false
}

func sameScalarJSON(first, second json.RawMessage) bool {
	left, err := canonicalJSON(first)
	if err != nil {
		return false
	}
	right, err := canonicalJSON(second)
	if err != nil {
		return false
	}
	return bytes.Equal(left, right)
}

func rejectSensitiveURL(u *url.URL) error {
	for key := range u.Query() {
		if sensitiveKey(key) {
			return fmt.Errorf("provenance endpoint contains prohibited sensitive query field %q", key)
		}
	}
	return nil
}

func isSHA256(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil && value == strings.ToLower(value)
}

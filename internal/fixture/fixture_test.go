package fixture

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const validLedgerKey = "AQIDBA=="

func tooManyLedgerKeysParams() json.RawMessage {
	keys := make([]string, MaxLedgerKeys+1)
	for i := range keys {
		keys[i] = validLedgerKey
	}
	params, err := json.Marshal(map[string]any{"keys": keys})
	if err != nil {
		panic(err)
	}
	return params
}

func validFixture(params string) Fixture {
	return Fixture{
		SchemaVersion: SchemaVersion,
		CapturedAt:    "2026-10-08T12:00:00Z",
		Provenance: Provenance{
			Endpoint:    "https://soroban-testnet.stellar.org",
			Network:     "testnet",
			ToolVersion: "dev",
		},
		Request: Request{
			JSONRPC: "2.0",
			ID:      json.RawMessage(`1`),
			Method:  "getLedgerEntries",
			Params:  json.RawMessage(params),
		},
		Response: Response{
			JSONRPC: "2.0",
			ID:      json.RawMessage(`1`),
			Result:  json.RawMessage(`{"latestLedger":10,"entries":[]}`),
		},
	}
}

func TestSealAndValidate(t *testing.T) {
	f := validFixture(`{"keys":["` + validLedgerKey + `"],"xdrFormat":"base64"}`)
	if err := Seal(&f); err != nil {
		t.Fatal(err)
	}
	if err := Validate(f); err != nil {
		t.Fatal(err)
	}
	if len(f.Integrity.SHA256) != 64 {
		t.Fatalf("unexpected hash %q", f.Integrity.SHA256)
	}
}

func TestCanonicalHashIgnoresObjectOrder(t *testing.T) {
	first := validFixture(`{"xdrFormat":"base64","keys":["` + validLedgerKey + `"]}`)
	second := validFixture(`{"keys":["` + validLedgerKey + `"],"xdrFormat":"base64"}`)
	if err := Seal(&first); err != nil {
		t.Fatal(err)
	}
	hash, err := Hash(second)
	if err != nil {
		t.Fatal(err)
	}
	if hash != first.Integrity.SHA256 {
		t.Fatalf("hash changed with object order: %s != %s", hash, first.Integrity.SHA256)
	}
}

func TestModifiedFixtureFailsIntegrity(t *testing.T) {
	f := validFixture(`{"keys":["` + validLedgerKey + `"]}`)
	if err := Seal(&f); err != nil {
		t.Fatal(err)
	}
	f.Response.Result = json.RawMessage(`{"latestLedger":11,"entries":[]}`)
	if err := Validate(f); err == nil || !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("expected integrity mismatch, got %v", err)
	}
}

func TestRejectsInvalidFixtures(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Fixture)
		want   string
	}{
		{"unknown method", func(f *Fixture) { f.Request.Method = "sendTransaction" }, "unsupported RPC method"},
		{"positional params", func(f *Fixture) { f.Request.Method = "getHealth"; f.Request.Params = json.RawMessage(`[]`) }, "empty object"},
		{"missing ledger keys", func(f *Fixture) { f.Request.Params = json.RawMessage(`{}`) }, "keys are required"},
		{"too many ledger keys", func(f *Fixture) { f.Request.Params = tooManyLedgerKeysParams() }, "at most"},
		{"bad xdr format", func(f *Fixture) { f.Request.Params = json.RawMessage(`{"keys":["AQIDBA=="],"xdrFormat":"xml"}`) }, "xdrFormat"},
		{"secret field", func(f *Fixture) { f.Request.Params = json.RawMessage(`{"keys":["AQIDBA=="],"api_key":"secret"}`) }, "sensitive"},
		{"token field", func(f *Fixture) { f.Request.Params = json.RawMessage(`{"keys":["AQIDBA=="],"X-Token":"secret"}`) }, "sensitive"},
		{"bad endpoint", func(f *Fixture) { f.Provenance.Endpoint = "http://localhost:8000" }, "HTTPS"},
		{"two response branches", func(f *Fixture) { f.Response.Error = &RPCError{Code: -1, Message: "error"} }, "exactly one"},
		{"mismatched response ID", func(f *Fixture) { f.Response.ID = json.RawMessage(`"1"`) }, "IDs must match"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := validFixture(`{"keys":["AQIDBA=="]}`)
			test.mutate(&f)
			if err := validateStructure(f); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("expected %q, got %v", test.want, err)
			}
		})
	}
}

func TestValidatePreservesJSONRPCIDTypes(t *testing.T) {
	f := validFixture(`{"keys":["AQIDBA=="]}`)
	for _, responseID := range []string{`"1"`, `null`, `1.0`} {
		t.Run(responseID, func(t *testing.T) {
			candidate := f
			candidate.Response.ID = json.RawMessage(responseID)
			if err := validateStructure(candidate); err == nil || !strings.Contains(err.Error(), "IDs must match") {
				t.Fatalf("response ID %s was accepted: %v", responseID, err)
			}
		})
	}
}

func TestSaveLoadAndRejectUnknownFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixture.json")
	f := validFixture(`{"keys":["AQIDBA=="]}`)
	if err := Save(path, f); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(loaded); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "\"integrity\"") {
		t.Fatal("saved fixture has no integrity field")
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatal(err)
	}
	object["unexpected"] = json.RawMessage(`true`)
	modified, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, modified, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("unknown top-level field was accepted")
	}
}

func TestZeroParameterMethodsAcceptOmittedAndEmptyParams(t *testing.T) {
	for _, params := range []json.RawMessage{nil, json.RawMessage(`{}`)} {
		f := validFixture(`{"keys":["AQIDBA=="]}`)
		f.Request.Method = "getHealth"
		f.Request.Params = params
		if err := validateStructure(f); err != nil {
			t.Fatalf("params %s: %v", params, err)
		}
	}
}

func TestSaveCreatesNestedPrivateFixtureAndLoadRejectsOversizedFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "fixture.json")
	if err := Save(path, validFixture(`{"keys":["AQIDBA=="]}`)); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Integrity == nil {
		t.Fatal("loaded fixture has no integrity")
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Fatalf("fixture permissions are %o, want 600", info.Mode().Perm())
		}
	}

	oversized := filepath.Join(dir, "oversized.json")
	if err := os.WriteFile(oversized, make([]byte, MaxFixtureBytes+1), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(oversized); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("oversized fixture was accepted: %v", err)
	}
}

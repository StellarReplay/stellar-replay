package replay

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
)

func TestReplayExactMatchNormalizesObjectsAndEchoesID(t *testing.T) {
	engine := mustEngine(t, validFixture("getLedgerEntries", `{"xdrFormat":"base64","keys":["AQIDBA=="]}`, `{"entries":[{"key":"one"}]}`))
	response, err := engine.Replay(fixture.Request{
		JSONRPC: "2.0", ID: json.RawMessage(`99`), Method: "getLedgerEntries",
		Params: json.RawMessage(`{"keys":["AQIDBA=="],"xdrFormat":"base64"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(response.ID) != "99" || string(response.Result) != `{"entries":[{"key":"one"}]}` {
		t.Fatalf("unexpected replay response: %+v", response)
	}
}

func TestReplayZeroParameterOmittedAndEmptyParamsMatch(t *testing.T) {
	engine := mustEngine(t, validFixture("getHealth", `{}`, `{"status":"healthy"}`))
	for _, id := range []string{"1", `"client-id"`} {
		response, err := engine.Replay(fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(id), Method: "getHealth"})
		if err != nil {
			t.Fatal(err)
		}
		if string(response.ID) != id {
			t.Fatalf("expected echoed id %s, got %s", id, response.ID)
		}
		response, err = engine.Replay(fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(id), Method: "getHealth", Params: json.RawMessage(`{}`)})
		if err != nil || string(response.Result) != `{"status":"healthy"}` {
			t.Fatalf("empty params did not match: response=%+v err=%v", response, err)
		}
	}
}

func TestReplayPreservesArrayOrder(t *testing.T) {
	engine := mustEngine(t, validFixture("getLedgerEntries", `{"keys":["AQIDBA==","BAUGBw=="]}`, `{"entries":[]}`))
	_, err := engine.Replay(fixture.Request{
		JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: "getLedgerEntries",
		Params: json.RawMessage(`{"keys":["BAUGBw==","AQIDBA=="]}`),
	})
	var replayErr *Error
	if !errors.As(err, &replayErr) || replayErr.Code != CodeInvalidParams || !strings.Contains(err.Error(), "no replay fixture matched") {
		t.Fatalf("array order was normalized unexpectedly: %v", err)
	}
}

func TestReplayRejectsMissesMethodAndInvalidRequests(t *testing.T) {
	engine := mustEngine(t, validFixture("getLedgerEntries", `{"keys":["AQIDBA=="]}`, `{"entries":[]}`))
	tests := []struct {
		name    string
		request fixture.Request
		code    int
		text    string
	}{
		{"parameter miss", fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: "getLedgerEntries", Params: json.RawMessage(`{"keys":["BAUGBw=="]}`)}, CodeInvalidParams, "no replay fixture matched"},
		{"type mismatch", fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: "getLedgerEntries", Params: json.RawMessage(`{"keys":[1]}`)}, CodeInvalidParams, "invalid parameters"},
		{"unknown method", fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: "sendTransaction"}, CodeMethodNotFound, "not supported"},
		{"bad envelope", fixture.Request{JSONRPC: "1.0", ID: json.RawMessage(`1`), Method: "getLedgerEntries", Params: json.RawMessage(`{"keys":["AQIDBA=="]}`)}, CodeInvalidRequest, "requires jsonrpc"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := engine.Replay(test.request)
			var replayErr *Error
			if !errors.As(err, &replayErr) || replayErr.Code != test.code || !strings.Contains(err.Error(), test.text) {
				t.Fatalf("expected code %d containing %q, got %v", test.code, test.text, err)
			}
		})
	}
}

func TestReplayJSONRejectsMalformedBatchAndUnknownFields(t *testing.T) {
	engine := mustEngine(t, validFixture("getHealth", `{}`, `{"status":"healthy"}`))
	tests := []string{
		`{"jsonrpc":"2.0","id":1,"method":"getHealth"} trailing`,
		`[{"jsonrpc":"2.0","id":1,"method":"getHealth"}]`,
		`{"jsonrpc":"2.0","id":1,"method":"getHealth","unexpected":true}`,
	}
	for _, raw := range tests {
		_, err := engine.ReplayJSON([]byte(raw))
		var replayErr *Error
		if !errors.As(err, &replayErr) || replayErr.Code != CodeInvalidRequest {
			t.Fatalf("expected invalid request for %s, got %v", raw, err)
		}
	}
}

func TestNewRejectsInvalidTamperedAndDuplicateFixtures(t *testing.T) {
	valid := validFixture("getHealth", `{}`, `{"status":"healthy"}`)
	if _, err := New([]fixture.Fixture{valid, valid}); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate fixture accepted: %v", err)
	}
	tampered := valid
	tampered.Response.Result = json.RawMessage(`{"status":"tampered"}`)
	if _, err := New([]fixture.Fixture{tampered}); err == nil || !strings.Contains(err.Error(), "integrity") {
		t.Fatalf("tampered fixture accepted: %v", err)
	}
}

func TestLoadAndErrorResponse(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "health.json")
	if err := fixture.Save(path, validFixture("getHealth", `{}`, `{"status":"healthy"}`)); err != nil {
		t.Fatal(err)
	}
	engine, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	response, err := engine.ReplayJSON([]byte(`{"jsonrpc":"2.0","id":"abc","method":"getHealth","params":{}}`))
	if err != nil || string(response.ID) != `"abc"` {
		t.Fatalf("loaded engine replay failed: response=%+v err=%v", response, err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	errorResponse := ErrorResponse(json.RawMessage(`7`), &Error{Code: CodeInvalidParams, Message: "safe miss"})
	if errorResponse.JSONRPC != "2.0" || string(errorResponse.ID) != "7" || errorResponse.Error.Code != CodeInvalidParams {
		t.Fatalf("unexpected error response: %+v", errorResponse)
	}
}

func TestReplayConcurrentReads(t *testing.T) {
	engine := mustEngine(t, validFixture("getHealth", `{}`, `{"status":"healthy"}`))
	const workers = 64
	errorsCh := make(chan error, workers)
	var wait sync.WaitGroup
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			id := strconv.Itoa(i)
			response, err := engine.Replay(fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(id), Method: "getHealth", Params: json.RawMessage(`{}`)})
			if err != nil {
				errorsCh <- err
				return
			}
			if string(response.ID) != id {
				errorsCh <- errors.New("response id was not echoed")
			}
		}(i)
	}
	wait.Wait()
	close(errorsCh)
	for err := range errorsCh {
		t.Fatal(err)
	}
}

func mustEngine(t *testing.T, candidates ...fixture.Fixture) *Engine {
	t.Helper()
	engine, err := New(candidates)
	if err != nil {
		t.Fatal(err)
	}
	return engine
}

func validFixture(method, params, result string) fixture.Fixture {
	candidate := fixture.Fixture{
		SchemaVersion: fixture.SchemaVersion,
		CapturedAt:    "2026-10-08T12:00:00Z",
		Provenance: fixture.Provenance{
			Endpoint: "https://soroban-testnet.stellar.org", Network: "testnet", ToolVersion: "test",
		},
		Request:  fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: method, Params: json.RawMessage(params)},
		Response: fixture.Response{JSONRPC: "2.0", ID: json.RawMessage(`1`), Result: json.RawMessage(result)},
	}
	if err := fixture.Seal(&candidate); err != nil {
		panic(err)
	}
	return candidate
}

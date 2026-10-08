package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
)

func TestRunHelpVersionAndUsageCodes(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"help"}, &stdout, &stderr); code != ExitOK || !strings.Contains(stdout.String(), "record") {
		t.Fatalf("help failed: code=%d out=%q err=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	if code := Run([]string{"version"}, &stdout, &stderr); code != ExitOK || stdout.String() != Version+"\n" {
		t.Fatalf("version failed: code=%d out=%q", code, stdout.String())
	}
	stdout.Reset()
	if code := Run([]string{"unknown"}, &stdout, &stderr); code != ExitUsage || !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("unknown command failed: code=%d err=%q", code, stderr.String())
	}
}

func TestRunValidateInspectAndReplayEndToEnd(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "health.json")
	candidate := cliFixture()
	if err := fixture.Save(path, candidate); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := Run([]string{"validate", "--fixture", path}, &stdout, &stderr); code != ExitOK || !strings.Contains(stdout.String(), "valid ") || stderr.Len() != 0 {
		t.Fatalf("validate failed: code=%d out=%q err=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	if code := Run([]string{"inspect", "--fixture", path}, &stdout, &stderr); code != ExitOK || !strings.Contains(stdout.String(), "method=getHealth") || !strings.Contains(stdout.String(), "response=result") {
		t.Fatalf("inspect failed: code=%d out=%q err=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	if code := Run([]string{"replay", "--fixture", path, "--method", "getHealth", "--params", "{}", "--id", `"client"`}, &stdout, &stderr); code != ExitOK {
		t.Fatalf("replay failed: code=%d out=%q err=%q", code, stdout.String(), stderr.String())
	}
	var response fixture.Response
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if string(response.ID) != `"client"` || string(response.Result) != `{"status":"healthy"}` {
		t.Fatalf("unexpected replay output: %+v", response)
	}

	requestPath := filepath.Join(dir, "request.json")
	if err := os.WriteFile(requestPath, []byte(`{"jsonrpc":"2.0","id":8,"method":"getHealth"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if code := Run([]string{"replay", "--fixture", path, "--request-file", requestPath}, &stdout, &stderr); code != ExitOK || !strings.Contains(stdout.String(), `"id":8`) {
		t.Fatalf("request-file replay failed: code=%d out=%q err=%q", code, stdout.String(), stderr.String())
	}
}

func TestRunRejectsBadArgumentsAndMissingFixtures(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := Run([]string{"validate"}, &stdout, &stderr); code != ExitUsage || !strings.Contains(stderr.String(), "requires --fixture") {
		t.Fatalf("validate usage failed: code=%d err=%q", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := Run([]string{"record", "--rpc-url", "http://localhost", "--method", "getHealth", "--network", "testnet", "--output", "out.json"}, &stdout, &stderr); code != ExitFailure || !strings.Contains(stderr.String(), "HTTPS") {
		t.Fatalf("record policy failure failed: code=%d err=%q", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := Run([]string{"replay", "--fixture", filepath.Join(t.TempDir(), "missing.json"), "--method", "getHealth"}, &stdout, &stderr); code != ExitFailure || !strings.Contains(stderr.String(), "load fixture") {
		t.Fatalf("missing fixture failed: code=%d err=%q", code, stderr.String())
	}
}

func cliFixture() fixture.Fixture {
	candidate := fixture.Fixture{
		SchemaVersion: fixture.SchemaVersion,
		CapturedAt:    "2026-10-08T12:00:00Z",
		Provenance: fixture.Provenance{
			Endpoint: "https://soroban-testnet.stellar.org", Network: "testnet", ToolVersion: "test",
		},
		Request:  fixture.Request{JSONRPC: "2.0", ID: json.RawMessage(`1`), Method: "getHealth", Params: json.RawMessage(`{}`)},
		Response: fixture.Response{JSONRPC: "2.0", ID: json.RawMessage(`1`), Result: json.RawMessage(`{"status":"healthy"}`)},
	}
	if err := fixture.Seal(&candidate); err != nil {
		panic(err)
	}
	return candidate
}

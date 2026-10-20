package fixture

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckedInEdgeCaseFixtures(t *testing.T) {
	retention := filepath.Join("..", "..", "fixtures", "edge-cases", "get-ledger-entries-retention.json")
	candidate, err := Load(retention)
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Response.Error == nil || candidate.Response.Error.Code != -32001 {
		t.Fatalf("retention fixture did not preserve its JSON-RPC error: %+v", candidate.Response)
	}
	if !strings.Contains(candidate.Response.Error.Message, "retention") {
		t.Fatalf("retention fixture message was lost: %q", candidate.Response.Error.Message)
	}
}

func TestInvalidEdgeCaseFixturesAreRejected(t *testing.T) {
	cases := []struct {
		name string
		file string
		want string
	}{
		{"malformed response JSON", "malformed-response.json", "decode fixture"},
		{"missing result or error", "response-missing-branch.json", "exactly one"},
		{"both result and error", "response-both-branches.json", "exactly one"},
		{"unsupported pagination cursor", "pagination-cursor.json", "unsupported parameter"},
		{"response ID mismatch", "response-id-mismatch.json", "IDs must match"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			path := filepath.Join("..", "..", "tests", "fixtures", "invalid", testCase.file)
			if _, err := Load(path); err == nil || !strings.Contains(err.Error(), testCase.want) {
				t.Fatalf("expected %q from %s, got %v", testCase.want, path, err)
			}
		})
	}
}

func TestRetentionErrorRemainsAValidResponseBranch(t *testing.T) {
	candidate := validFixture(`{"keys":["AQIDBA=="]}`)
	candidate.Response.Result = nil
	candidate.Response.Error = &RPCError{Code: -32001, Message: "ledger is outside provider retention window"}
	if err := Seal(&candidate); err != nil {
		t.Fatal(err)
	}
	if err := Validate(candidate); err != nil {
		t.Fatal(err)
	}
	if candidate.Response.Error == nil {
		t.Fatal("retention error branch was not retained")
	}
}

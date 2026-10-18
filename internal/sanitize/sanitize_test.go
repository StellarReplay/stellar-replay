package sanitize

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
)

func TestResponseRedactsSensitiveKeysAndSecretValues(t *testing.T) {
	response := fixture.Response{
		JSONRPC: "2.0", ID: json.RawMessage(`1`),
		Result: json.RawMessage(`{"account":"G...safe","api_key":"abc","nested":{"password":"pw"},"pem":"-----BEGIN PRIVATE KEY-----secret-----END PRIVATE KEY-----","items":[{"authorization":"Bearer abc"}]}`),
	}
	if err := Response(&response); err != nil {
		t.Fatal(err)
	}
	var value map[string]any
	if err := json.Unmarshal(response.Result, &value); err != nil {
		t.Fatal(err)
	}
	if value["account"] != "G...safe" || value["api_key"] != redactedValue || value["pem"] != redactedValue {
		t.Fatalf("unexpected redaction: %#v", value)
	}
	nested := value["nested"].(map[string]any)
	if nested["password"] != redactedValue {
		t.Fatalf("nested secret was not redacted: %#v", nested)
	}
	items := value["items"].([]any)
	if items[0].(map[string]any)["authorization"] != redactedValue {
		t.Fatalf("array secret was not redacted: %#v", items)
	}
}

func FuzzResponseNoPanic(f *testing.F) {
	f.Add([]byte(`{"safe":"value","api_key":"secret"}`))
	f.Add([]byte(`[{"authorization":"Bearer secret"}]`))
	f.Add([]byte(`null`))
	f.Fuzz(func(t *testing.T, result []byte) {
		response := &fixture.Response{Result: result}
		_ = Response(response)
	})
}

func TestResponseRejectsMalformedJSON(t *testing.T) {
	response := fixture.Response{Result: json.RawMessage(`{"ok":true} trailing`)}
	if err := Response(&response); err == nil || !strings.Contains(err.Error(), "trailing") {
		t.Fatalf("expected malformed response error, got %v", err)
	}
}

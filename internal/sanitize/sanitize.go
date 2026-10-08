// Package sanitize applies conservative, explicit redaction to captured JSON.
package sanitize

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
)

const redactedValue = "[REDACTED]"

// Response redacts sensitive response fields and obvious secret material before
// a response is sealed into a fixture. It preserves object shape and array order.
// Unknown secret formats are not claimed to be detectable automatically.
func Response(response *fixture.Response) error {
	if response == nil {
		return errors.New("response is nil")
	}
	var err error
	response.Result, err = redactJSON(response.Result, "response.result")
	if err != nil {
		return err
	}
	if response.Error != nil {
		response.Error.Data, err = redactJSON(response.Error.Data, "response.error.data")
		if err != nil {
			return err
		}
	}
	return nil
}

func redactJSON(raw json.RawMessage, location string) (json.RawMessage, error) {
	if len(raw) == 0 {
		return raw, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("%s is invalid JSON: %w", location, err)
	}
	if err := ensureEOF(decoder); err != nil {
		return nil, fmt.Errorf("%s: %w", location, err)
	}
	redacted := redactValue(value)
	encoded, err := json.Marshal(redacted)
	if err != nil {
		return nil, fmt.Errorf("%s cannot be encoded: %w", location, err)
	}
	return encoded, nil
}

func redactValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			if sensitiveKey(key) {
				typed[key] = redactedValue
				continue
			}
			typed[key] = redactValue(child)
		}
	case []any:
		for i, child := range typed {
			typed[i] = redactValue(child)
		}
	case string:
		if sensitiveValue(typed) {
			return redactedValue
		}
	}
	return value
}

func sensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("_", "", "-", "", " ", "").Replace(key))
	for _, prohibited := range []string{
		"privatekey", "secret", "seedphrase", "mnemonic", "password", "authorization",
		"accesstoken", "refreshtoken", "apikey", "bearer", "signingkey", "signingmaterial",
	} {
		if strings.Contains(normalized, prohibited) {
			return true
		}
	}
	return false
}

func sensitiveValue(value string) bool {
	upper := strings.ToUpper(value)
	return strings.Contains(upper, "-----BEGIN ") && strings.Contains(upper, "PRIVATE KEY-----") ||
		strings.HasPrefix(strings.ToLower(value), "bearer ")
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

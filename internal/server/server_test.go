package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
	"github.com/stellar-replay/stellar-replay/internal/replay"
)

func TestServerStartsServesAndShutsDown(t *testing.T) {
	server := mustServer(t, replayFixture("getHealth", `{}`, `{"status":"healthy"}`), Config{})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Shutdown(context.Background()) }()
	if !strings.HasPrefix(server.Endpoint(), "http://127.0.0.1:") {
		t.Fatalf("unexpected endpoint: %s", server.Endpoint())
	}

	response, err := http.Post(server.Endpoint(), "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":42,"method":"getHealth"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK || response.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected HTTP response: %s %s", response.Status, response.Header.Get("Content-Type"))
	}
	var envelope fixture.Response
	if err := json.NewDecoder(response.Body).Decode(&envelope); err != nil {
		t.Fatal(err)
	}
	if string(envelope.ID) != "42" || string(envelope.Result) != `{"status":"healthy"}` {
		t.Fatalf("unexpected JSON-RPC response: %+v", envelope)
	}

	if err := server.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := server.Wait(); err != nil {
		t.Fatal(err)
	}
	_, err = http.Post(server.Endpoint(), "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"getHealth"}`))
	if err == nil {
		t.Fatal("request succeeded after shutdown")
	}
}

func TestServerRejectsMalformedUnsupportedMissAndBatch(t *testing.T) {
	server := mustServer(t, replayFixture("getHealth", `{}`, `{"status":"healthy"}`), Config{})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Shutdown(context.Background()) }()

	tests := []struct {
		name string
		body string
		code int
		id   string
	}{
		{"malformed", `{not-json`, replay.CodeInvalidRequest, "null"},
		{"unsupported", `{"jsonrpc":"2.0","id":"x","method":"sendTransaction"}`, replay.CodeMethodNotFound, `"x"`},
		{"miss", `{"jsonrpc":"2.0","id":9,"method":"getNetwork"}`, replay.CodeInvalidParams, "9"},
		{"batch", `[{"jsonrpc":"2.0","id":1,"method":"getHealth"}]`, replay.CodeInvalidRequest, "null"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := http.Post(server.Endpoint(), "application/json", strings.NewReader(test.body))
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			var envelope fixture.Response
			if err := json.NewDecoder(response.Body).Decode(&envelope); err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != http.StatusOK || envelope.Error == nil || envelope.Error.Code != test.code || string(envelope.ID) != test.id {
				t.Fatalf("unexpected error envelope: status=%d response=%+v", response.StatusCode, envelope)
			}
		})
	}
}

func TestServerRejectsTransportMisuseAndOversizedBody(t *testing.T) {
	server := mustServer(t, replayFixture("getHealth", `{}`, `{"status":"healthy"}`), Config{MaxRequestBytes: 64})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Shutdown(context.Background()) }()

	getResponse, err := http.Get(server.Endpoint())
	if err != nil {
		t.Fatal(err)
	}
	getResponse.Body.Close()
	if getResponse.StatusCode != http.StatusMethodNotAllowed || getResponse.Header.Get("Allow") != http.MethodPost {
		t.Fatalf("unexpected GET response: %s allow=%s", getResponse.Status, getResponse.Header.Get("Allow"))
	}
	pathResponse, err := http.Post(server.Endpoint()+"/rpc", "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"getHealth"}`))
	if err != nil {
		t.Fatal(err)
	}
	pathResponse.Body.Close()
	if pathResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected path response: %s", pathResponse.Status)
	}

	large := strings.Repeat("x", 65)
	largeResponse, err := http.Post(server.Endpoint(), "application/json", strings.NewReader(large))
	if err != nil {
		t.Fatal(err)
	}
	defer largeResponse.Body.Close()
	var envelope fixture.Response
	if err := json.NewDecoder(largeResponse.Body).Decode(&envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Error == nil || envelope.Error.Code != replay.CodeInvalidRequest || !strings.Contains(envelope.Error.Message, "exceeds") {
		t.Fatalf("unexpected oversized response: %+v", envelope)
	}
}

func TestServerConcurrentRequests(t *testing.T) {
	server := mustServer(t, replayFixture("getHealth", `{}`, `{"status":"healthy"}`), Config{})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Shutdown(context.Background()) }()

	const requests = 64
	errorsCh := make(chan error, requests)
	var wait sync.WaitGroup
	client := &http.Client{Timeout: 5 * time.Second}
	for i := 0; i < requests; i++ {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			body := `{"jsonrpc":"2.0","id":` + string(rune('0'+i%10)) + `,"method":"getHealth","params":{}}`
			response, err := client.Post(server.Endpoint(), "application/json", strings.NewReader(body))
			if err != nil {
				errorsCh <- err
				return
			}
			defer response.Body.Close()
			var envelope fixture.Response
			if err := json.NewDecoder(response.Body).Decode(&envelope); err != nil {
				errorsCh <- err
				return
			}
			if envelope.Error != nil || string(envelope.Result) != `{"status":"healthy"}` {
				errorsCh <- errors.New("concurrent replay returned an unexpected response")
			}
		}(i)
	}
	wait.Wait()
	close(errorsCh)
	for err := range errorsCh {
		t.Fatal(err)
	}
}

func TestServerConfigurationIsLoopbackOnly(t *testing.T) {
	engine := mustEngine(t, replayFixture("getHealth", `{}`, `{"status":"healthy"}`))
	for _, addr := range []string{"0.0.0.0:0", ":0", "192.168.1.10:8787"} {
		if _, err := New(engine, Config{Addr: addr}); err == nil {
			t.Fatalf("non-loopback address accepted: %s", addr)
		}
	}
	if _, err := New(engine, Config{Addr: "127.0.0.1:0"}); err != nil {
		t.Fatal(err)
	}
}

func mustServer(t *testing.T, candidate fixture.Fixture, cfg Config) *Server {
	t.Helper()
	engine := mustEngine(t, candidate)
	server, err := New(engine, cfg)
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func mustEngine(t *testing.T, candidate fixture.Fixture) *replay.Engine {
	t.Helper()
	engine, err := replay.New([]fixture.Fixture{candidate})
	if err != nil {
		t.Fatal(err)
	}
	return engine
}

func replayFixture(method, params, result string) fixture.Fixture {
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

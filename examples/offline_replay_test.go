package examples

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
	"github.com/stellar-replay/stellar-replay/internal/replay"
	"github.com/stellar-replay/stellar-replay/internal/server"
)

// TestOfflineReplayWorkflow is the contributor-facing example for the normal
// test path: checked-in fixtures, loopback server, repeated replay, and clients
// guarded against all non-loopback network access.
func TestOfflineReplayWorkflow(t *testing.T) {
	paths := []string{
		"../fixtures/get-health.json",
		"../fixtures/get-latest-ledger.json",
		"../fixtures/get-network.json",
		"../fixtures/get-ledger-entries.json",
	}
	engine, err := replay.Load(paths...)
	if err != nil {
		t.Fatal(err)
	}
	localServer, err := server.New(engine, server.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := localServer.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = localServer.Shutdown(context.Background()) }()

	client := &http.Client{Transport: loopbackOnlyTransport{base: http.DefaultTransport}}
	requestBody := []byte(`{"jsonrpc":"2.0","id":7,"method":"getLatestLedger"}`)
	first := doReplayRequest(t, client, localServer.Endpoint(), requestBody)
	second := doReplayRequest(t, client, localServer.Endpoint(), requestBody)
	if !bytes.Equal(first, second) {
		t.Fatalf("repeated offline replay was not deterministic:\n%s\n%s", first, second)
	}
	if !bytes.Contains(first, []byte(`"sequence":123456`)) {
		t.Fatalf("fixture response was not served: %s", first)
	}

	const concurrent = 32
	errorsCh := make(chan error, concurrent)
	var wait sync.WaitGroup
	for i := 0; i < concurrent; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			body, err := requestReplay(client, localServer.Endpoint(), requestBody)
			if err != nil {
				errorsCh <- err
				return
			}
			if !bytes.Equal(body, first) {
				errorsCh <- errors.New("concurrent replay response changed")
			}
		}()
	}
	wait.Wait()
	close(errorsCh)
	for err := range errorsCh {
		t.Fatal(err)
	}
}

func TestCheckedInFixtureMatrix(t *testing.T) {
	paths := []string{
		"../fixtures/get-health.json",
		"../fixtures/get-latest-ledger.json",
		"../fixtures/get-network.json",
		"../fixtures/get-ledger-entries.json",
	}
	for _, path := range paths {
		candidate, err := fixture.Load(path)
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		if candidate.Provenance.Network != "testnet" || candidate.Integrity == nil {
			t.Fatalf("%s is missing safe provenance or integrity", path)
		}
	}
}

func doReplayRequest(t *testing.T, client *http.Client, endpoint string, body []byte) []byte {
	t.Helper()
	result, err := requestReplay(client, endpoint, body)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func requestReplay(client *http.Client, endpoint string, body []byte) ([]byte, error) {
	response, err := client.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected replay status: %s", response.Status)
	}
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type loopbackOnlyTransport struct {
	base http.RoundTripper
}

func (transport loopbackOnlyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	host := request.URL.Hostname()
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return nil, fmt.Errorf("network guard blocked non-loopback destination %q", host)
	}
	return transport.base.RoundTrip(request)
}

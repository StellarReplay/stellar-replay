// Package server exposes validated replay fixtures through a loopback HTTP
// JSON-RPC endpoint.
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
	"github.com/stellar-replay/stellar-replay/internal/replay"
)

const (
	DefaultAddr              = "127.0.0.1:0"
	DefaultMaxRequestBytes   = fixture.MaxFixtureBytes
	DefaultReadHeaderTimeout = 5 * time.Second
)

var (
	ErrAlreadyStarted = errors.New("replay server is already started")
	ErrNotStarted     = errors.New("replay server is not started")
)

// Config defines the local server boundary. Addr must remain loopback-only.
type Config struct {
	Addr              string
	MaxRequestBytes   int64
	ReadHeaderTimeout time.Duration
}

// Server serves one immutable replay engine until Shutdown is called.
type Server struct {
	engine     *replay.Engine
	httpServer *http.Server
	addr       string

	mu       sync.Mutex
	listener net.Listener
	started  bool
	done     chan struct{}
	serveErr chan error
}

// New creates a loopback-only server. It does not open a listener until Start.
func New(engine *replay.Engine, cfg Config) (*Server, error) {
	if engine == nil {
		return nil, errors.New("replay engine is required")
	}
	addr := cfg.Addr
	if addr == "" {
		addr = DefaultAddr
	}
	if err := validateLoopbackAddr(addr); err != nil {
		return nil, err
	}
	maxBytes := cfg.MaxRequestBytes
	if maxBytes == 0 {
		maxBytes = DefaultMaxRequestBytes
	}
	if maxBytes <= 0 || maxBytes > fixture.MaxFixtureBytes {
		return nil, fmt.Errorf("max request bytes must be between 1 and %d", fixture.MaxFixtureBytes)
	}
	headerTimeout := cfg.ReadHeaderTimeout
	if headerTimeout == 0 {
		headerTimeout = DefaultReadHeaderTimeout
	}
	if headerTimeout < 0 {
		return nil, errors.New("read header timeout must be positive")
	}

	server := &Server{
		engine:   engine,
		addr:     addr,
		done:     make(chan struct{}),
		serveErr: make(chan error, 1),
	}
	server.httpServer = &http.Server{
		Handler:           server.handler(maxBytes),
		ReadHeaderTimeout: headerTimeout,
	}
	return server, nil
}

// Start binds the configured loopback address and begins serving in a goroutine.
func (server *Server) Start() error {
	if server == nil {
		return errors.New("replay server is nil")
	}
	server.mu.Lock()
	defer server.mu.Unlock()
	if server.started {
		return ErrAlreadyStarted
	}
	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", server.addr, err)
	}
	server.listener = listener
	server.started = true
	go server.serve(listener)
	return nil
}

// Addr returns the bound listener address after Start, or the configured address
// before Start.
func (server *Server) Addr() string {
	if server == nil {
		return ""
	}
	server.mu.Lock()
	defer server.mu.Unlock()
	if server.listener != nil {
		return server.listener.Addr().String()
	}
	return server.addr
}

// Endpoint returns an HTTP URL suitable for a local JSON-RPC client.
func (server *Server) Endpoint() string {
	return "http://" + server.Addr()
}

// Wait blocks until serving stops and returns any unexpected serve error.
func (server *Server) Wait() error {
	if server == nil {
		return errors.New("replay server is nil")
	}
	server.mu.Lock()
	started := server.started
	server.mu.Unlock()
	if !started {
		return ErrNotStarted
	}
	<-server.done
	select {
	case err := <-server.serveErr:
		return err
	default:
		return nil
	}
}

// Shutdown stops accepting requests and waits for active handlers to finish.
func (server *Server) Shutdown(ctx context.Context) error {
	if server == nil {
		return errors.New("replay server is nil")
	}
	server.mu.Lock()
	started := server.started
	server.mu.Unlock()
	if !started {
		return ErrNotStarted
	}
	if err := server.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown replay server: %w", err)
	}
	return nil
}

func (server *Server) serve(listener net.Listener) {
	err := server.httpServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		server.serveErr <- err
	}
	close(server.done)
}

func (server *Server) handler(maxBytes int64) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(responseWriter, request)
			return
		}
		if request.Method != http.MethodPost {
			responseWriter.Header().Set("Allow", http.MethodPost)
			http.Error(responseWriter, "method must be POST", http.StatusMethodNotAllowed)
			return
		}

		body, err := readRequestBody(responseWriter, request, maxBytes)
		if err != nil {
			writeRPCError(responseWriter, nil, &replay.Error{Code: replay.CodeInvalidRequest, Message: err.Error()})
			return
		}
		id := extractID(body)
		response, replayErr := server.engine.ReplayJSON(body)
		if replayErr != nil {
			writeRPCError(responseWriter, id, replayErr)
			return
		}
		writeJSON(responseWriter, response)
	})
}

func readRequestBody(responseWriter http.ResponseWriter, request *http.Request, maxBytes int64) ([]byte, error) {
	limited := http.MaxBytesReader(responseWriter, request.Body, maxBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), "request body too large") {
			return nil, fmt.Errorf("request body exceeds %d bytes", maxBytes)
		}
		return nil, fmt.Errorf("read request body: %w", err)
	}
	if int64(len(body)) > maxBytes {
		return nil, fmt.Errorf("request body exceeds %d bytes", maxBytes)
	}
	return body, nil
}

func writeRPCError(responseWriter http.ResponseWriter, id json.RawMessage, err error) {
	writeJSON(responseWriter, replay.ErrorResponse(id, err))
}

func writeJSON(responseWriter http.ResponseWriter, value any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(responseWriter).Encode(value)
}

func extractID(raw []byte) json.RawMessage {
	var object map[string]json.RawMessage
	decoder := json.NewDecoder(bytes.NewReader(raw))
	if decoder.Decode(&object) != nil {
		return nil
	}
	id, ok := object["id"]
	if !ok || !scalarJSON(id) {
		return nil
	}
	return append(json.RawMessage(nil), id...)
}

func validateLoopbackAddr(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("server address must be host:port: %w", err)
	}
	if host == "localhost" {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return errors.New("server address must be loopback-only")
	}
	return nil
}

func scalarJSON(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if decoder.Decode(&value) != nil {
		return false
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		return false
	}
	switch value.(type) {
	case nil, string, json.Number:
		return true
	default:
		return false
	}
}

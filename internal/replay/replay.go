// Package replay matches validated fixtures without contacting a live endpoint.
package replay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/stellar-replay/stellar-replay/internal/fixture"
)

const (
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603
)

// Error is a deterministic replay error suitable for later JSON-RPC mapping.
type Error struct {
	Code    int
	Message string
	Data    json.RawMessage
}

func (e *Error) Error() string {
	if e == nil {
		return "replay error"
	}
	return fmt.Sprintf("replay error %d: %s", e.Code, e.Message)
}

// Engine is an immutable, concurrency-safe index of validated fixtures.
type Engine struct {
	mu    sync.RWMutex
	byKey map[string]fixture.Fixture
}

// New validates and indexes fixtures. Duplicate canonical request keys fail
// closed because choosing between them would make replay nondeterministic.
func New(fixtures []fixture.Fixture) (*Engine, error) {
	engine := &Engine{
		byKey: make(map[string]fixture.Fixture, len(fixtures)),
	}
	for i, candidate := range fixtures {
		if err := fixture.Validate(candidate); err != nil {
			return nil, fmt.Errorf("fixture %d is invalid: %w", i, err)
		}
		key, err := requestKey(candidate.Request)
		if err != nil {
			return nil, fmt.Errorf("fixture %d request is invalid: %w", i, err)
		}
		if _, exists := engine.byKey[key]; exists {
			return nil, fmt.Errorf("duplicate fixture request key for %q", candidate.Request.Method)
		}
		engine.byKey[key] = candidate
	}
	return engine, nil
}

// Load reads and validates fixture files before indexing them.
func Load(paths ...string) (*Engine, error) {
	fixtures := make([]fixture.Fixture, 0, len(paths))
	for _, path := range paths {
		candidate, err := fixture.Load(path)
		if err != nil {
			return nil, fmt.Errorf("load fixture %q: %w", path, err)
		}
		fixtures = append(fixtures, candidate)
	}
	return New(fixtures)
}

// Replay returns the stored response for an exact canonical request match.
// The returned response has the incoming request id and no other mutation.
func (engine *Engine) Replay(request fixture.Request) (fixture.Response, error) {
	if engine == nil {
		return fixture.Response{}, &Error{Code: CodeInternalError, Message: "replay engine is nil"}
	}
	if err := validateRequestEnvelope(request); err != nil {
		return fixture.Response{}, err
	}
	if !fixture.Supported(request.Method) {
		return fixture.Response{}, &Error{Code: CodeMethodNotFound, Message: fmt.Sprintf("method %q is not supported", request.Method)}
	}
	if err := fixture.ValidateRequest(request); err != nil {
		return fixture.Response{}, &Error{Code: CodeInvalidParams, Message: "invalid parameters"}
	}
	key, err := requestKey(request)
	if err != nil {
		return fixture.Response{}, &Error{Code: CodeInvalidParams, Message: "invalid parameters"}
	}

	engine.mu.RLock()
	stored, ok := engine.byKey[key]
	engine.mu.RUnlock()
	if !ok {
		return fixture.Response{}, &Error{
			Code:    CodeInvalidParams,
			Message: fmt.Sprintf("no replay fixture matched method %q and its canonical parameters", request.Method),
		}
	}
	response := cloneResponse(stored.Response)
	response.ID = cloneRaw(request.ID)
	return response, nil
}

// ReplayJSON parses one non-batch JSON-RPC request and replays it.
func (engine *Engine) ReplayJSON(raw []byte) (fixture.Response, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var request fixture.Request
	if err := decoder.Decode(&request); err != nil {
		return fixture.Response{}, &Error{Code: CodeInvalidRequest, Message: "request is not valid JSON-RPC"}
	}
	if err := ensureEOF(decoder); err != nil {
		return fixture.Response{}, &Error{Code: CodeInvalidRequest, Message: "request must contain one JSON value"}
	}
	return engine.Replay(request)
}

// ErrorResponse converts a replay error into a JSON-RPC error envelope.
func ErrorResponse(id json.RawMessage, err error) fixture.Response {
	if err == nil {
		return fixture.Response{JSONRPC: "2.0", ID: cloneRaw(id), Error: &fixture.RPCError{Code: CodeInternalError, Message: "unknown replay error"}}
	}
	var replayErr *Error
	if !errors.As(err, &replayErr) {
		replayErr = &Error{Code: CodeInternalError, Message: "replay failed"}
	}
	return fixture.Response{
		JSONRPC: "2.0",
		ID:      cloneRaw(id),
		Error:   &fixture.RPCError{Code: replayErr.Code, Message: replayErr.Message, Data: cloneRaw(replayErr.Data)},
	}
}

func requestKey(request fixture.Request) (string, error) {
	params := request.Params
	if request.Method == "getHealth" || request.Method == "getLatestLedger" || request.Method == "getNetwork" {
		if len(params) == 0 {
			params = json.RawMessage(`{}`)
		}
	}
	canonical, err := fixture.CanonicalJSONValue(params)
	if err != nil {
		return "", err
	}
	return request.Method + "\x00" + string(canonical), nil
}

func validateRequestEnvelope(request fixture.Request) error {
	if request.JSONRPC != "2.0" || !scalarJSON(request.ID) || request.Method == "" {
		return &Error{Code: CodeInvalidRequest, Message: "request requires jsonrpc 2.0, scalar id, and method"}
	}
	return nil
}

func cloneResponse(response fixture.Response) fixture.Response {
	cloned := response
	cloned.ID = cloneRaw(response.ID)
	cloned.Result = cloneRaw(response.Result)
	if response.Error != nil {
		cloned.Error = &fixture.RPCError{
			Code: response.Error.Code, Message: response.Error.Message, Data: cloneRaw(response.Error.Data),
		}
	}
	return cloned
}

func cloneRaw(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	return append(json.RawMessage(nil), raw...)
}

func scalarJSON(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
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

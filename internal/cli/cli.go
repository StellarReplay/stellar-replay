// Package cli implements the stellar-replay command-line interface.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/stellar-replay/stellar-replay/internal/capture"
	"github.com/stellar-replay/stellar-replay/internal/fixture"
	"github.com/stellar-replay/stellar-replay/internal/replay"
	"github.com/stellar-replay/stellar-replay/internal/server"
)

const (
	ExitOK      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

// Version is overridden by release builds through -ldflags. Development builds
// retain an explicit suffix so they are not confused with a published binary.
var Version = "0.1.0-dev"

// Run executes one CLI invocation and returns a stable process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writeUsage(stderr)
		return ExitUsage
	}
	switch args[0] {
	case "help", "--help", "-h":
		writeUsage(stdout)
		return ExitOK
	case "version", "--version":
		fmt.Fprintln(stdout, Version)
		return ExitOK
	case "record":
		return runRecord(args[1:], stdout, stderr)
	case "inspect":
		return runInspect(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "replay":
		return runReplay(args[1:], stdout, stderr)
	case "serve":
		return runServe(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		writeUsage(stderr)
		return ExitUsage
	}
}

func runRecord(args []string, stdout, stderr io.Writer) int {
	flags := newFlagSet("record", stderr)
	rpcURL := flags.String("rpc-url", "", "explicit HTTPS RPC endpoint")
	method := flags.String("method", "", "supported RPC method")
	params := flags.String("params", "", "JSON-RPC params JSON")
	paramsFile := flags.String("params-file", "", "file containing params JSON, or - for stdin")
	id := flags.String("id", "1", "scalar JSON-RPC request id")
	network := flags.String("network", "", "network label recorded in provenance")
	toolVersion := flags.String("tool-version", Version, "tool version recorded in provenance")
	output := flags.String("output", "", "fixture output path")
	timeout := flags.Duration("timeout", capture.DefaultTimeout, "capture timeout")
	if err := flags.Parse(args); err != nil {
		return ExitUsage
	}
	if *rpcURL == "" || *method == "" || *network == "" || *output == "" {
		return usageError(stderr, "record requires --rpc-url, --method, --network, and --output")
	}
	paramJSON, err := readJSONOption(*params, *paramsFile, "params")
	if err != nil {
		return operationalError(stderr, err)
	}
	idJSON, err := parseJSON(*id, "id")
	if err != nil {
		return usageError(stderr, err.Error())
	}
	captured, err := capture.Record(context.Background(), capture.Config{
		Endpoint: *rpcURL, Network: *network, ToolVersion: *toolVersion,
		Method: *method, ID: idJSON, Params: paramJSON, Timeout: *timeout,
	}, *output)
	if err != nil {
		return operationalError(stderr, err)
	}
	fmt.Fprintf(stdout, "recorded %s\nsha256=%s\n", *output, captured.Integrity.SHA256)
	return ExitOK
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	flags := newFlagSet("validate", stderr)
	path := flags.String("fixture", "", "fixture path")
	if err := flags.Parse(args); err != nil {
		return ExitUsage
	}
	if *path == "" {
		return usageError(stderr, "validate requires --fixture PATH")
	}
	candidate, err := fixture.Load(*path)
	if err != nil {
		return operationalError(stderr, err)
	}
	fmt.Fprintf(stdout, "valid %s\nsha256=%s\n", *path, candidate.Integrity.SHA256)
	return ExitOK
}

func runInspect(args []string, stdout, stderr io.Writer) int {
	flags := newFlagSet("inspect", stderr)
	path := flags.String("fixture", "", "fixture path")
	if err := flags.Parse(args); err != nil {
		return ExitUsage
	}
	if *path == "" {
		return usageError(stderr, "inspect requires --fixture PATH")
	}
	candidate, err := fixture.Load(*path)
	if err != nil {
		return operationalError(stderr, err)
	}
	responseKind := "error"
	if candidate.Response.Error == nil {
		responseKind = "result"
	}
	fmt.Fprintf(stdout, "schemaVersion=%s\nmethod=%s\nendpoint=%s\nnetwork=%s\ncapturedAt=%s\nresponse=%s\nsha256=%s\n",
		candidate.SchemaVersion, candidate.Request.Method, candidate.Provenance.Endpoint,
		candidate.Provenance.Network, candidate.CapturedAt, responseKind, candidate.Integrity.SHA256)
	return ExitOK
}

func runReplay(args []string, stdout, stderr io.Writer) int {
	flags := newFlagSet("replay", stderr)
	fixtures := stringList{}
	flags.Var(&fixtures, "fixture", "fixture path (repeatable)")
	requestFile := flags.String("request-file", "", "JSON-RPC request file, or - for stdin")
	method := flags.String("method", "", "supported RPC method")
	params := flags.String("params", "", "JSON-RPC params JSON")
	paramsFile := flags.String("params-file", "", "file containing params JSON, or - for stdin")
	id := flags.String("id", "1", "scalar JSON-RPC request id")
	if err := flags.Parse(args); err != nil {
		return ExitUsage
	}
	if len(fixtures) == 0 {
		return usageError(stderr, "replay requires at least one --fixture PATH")
	}
	engine, err := replay.Load(fixtures...)
	if err != nil {
		return operationalError(stderr, err)
	}
	var response fixture.Response
	if *requestFile != "" {
		raw, readErr := readInput(*requestFile)
		if readErr != nil {
			return operationalError(stderr, readErr)
		}
		response, err = engine.ReplayJSON(raw)
	} else {
		if *method == "" {
			return usageError(stderr, "replay requires --method or --request-file")
		}
		paramJSON, paramErr := readJSONOption(*params, *paramsFile, "params")
		if paramErr != nil {
			return operationalError(stderr, paramErr)
		}
		idJSON, idErr := parseJSON(*id, "id")
		if idErr != nil {
			return usageError(stderr, idErr.Error())
		}
		response, err = engine.Replay(fixture.Request{JSONRPC: "2.0", ID: idJSON, Method: *method, Params: paramJSON})
	}
	if err != nil {
		_ = json.NewEncoder(stderr).Encode(replay.ErrorResponse(nil, err))
		return ExitFailure
	}
	return writeJSON(stdout, response)
}

func runServe(args []string, stdout, stderr io.Writer) int {
	flags := newFlagSet("serve", stderr)
	fixtures := stringList{}
	flags.Var(&fixtures, "fixture", "fixture path (repeatable)")
	listen := flags.String("listen", "127.0.0.1:8787", "loopback listen address")
	if err := flags.Parse(args); err != nil {
		return ExitUsage
	}
	if len(fixtures) == 0 {
		return usageError(stderr, "serve requires at least one --fixture PATH")
	}
	engine, err := replay.Load(fixtures...)
	if err != nil {
		return operationalError(stderr, err)
	}
	localServer, err := server.New(engine, server.Config{Addr: *listen})
	if err != nil {
		return operationalError(stderr, err)
	}
	if err := localServer.Start(); err != nil {
		return operationalError(stderr, err)
	}
	fmt.Fprintf(stdout, "serving %s\n", localServer.Endpoint())

	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	wait := make(chan error, 1)
	go func() { wait <- localServer.Wait() }()
	select {
	case <-signalContext.Done():
		if err := localServer.Shutdown(context.Background()); err != nil {
			return operationalError(stderr, err)
		}
		return ExitOK
	case err := <-wait:
		if err != nil {
			return operationalError(stderr, err)
		}
		return ExitOK
	}
}

func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	return flags
}

func readJSONOption(value, path, name string) (json.RawMessage, error) {
	if value != "" && path != "" {
		return nil, fmt.Errorf("--%s and --%s-file cannot both be set", name, name)
	}
	if path != "" {
		raw, err := readInput(path)
		if err != nil {
			return nil, err
		}
		if !json.Valid(raw) {
			return nil, fmt.Errorf("%s file is not valid JSON", name)
		}
		return raw, nil
	}
	if value == "" {
		return nil, nil
	}
	return parseJSON(value, name)
}

func parseJSON(value, name string) (json.RawMessage, error) {
	raw := json.RawMessage(value)
	if !json.Valid(raw) {
		return nil, fmt.Errorf("%s must be valid JSON", name)
	}
	return raw, nil
}

func readInput(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func writeJSON(writer io.Writer, value any) int {
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		return ExitFailure
	}
	return ExitOK
}

func operationalError(stderr io.Writer, err error) int {
	fmt.Fprintf(stderr, "error: %v\n", err)
	return ExitFailure
}

func usageError(stderr io.Writer, message string) int {
	fmt.Fprintf(stderr, "usage error: %s\n", message)
	return ExitUsage
}

func writeUsage(writer io.Writer) {
	fmt.Fprintln(writer, "stellar-replay — deterministic Stellar RPC capture and replay")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "Commands:")
	fmt.Fprintln(writer, "  record   Capture one explicit HTTPS RPC interaction")
	fmt.Fprintln(writer, "  inspect  Print fixture metadata")
	fmt.Fprintln(writer, "  validate Validate fixture schema and integrity")
	fmt.Fprintln(writer, "  replay   Replay a fixture offline")
	fmt.Fprintln(writer, "  serve    Serve fixtures on a loopback HTTP endpoint")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "Use 'stellar-replay COMMAND -h' for command options.")
}

type stringList []string

func (list *stringList) String() string { return strings.Join(*list, ",") }

func (list *stringList) Set(value string) error {
	if value == "" {
		return errors.New("value cannot be empty")
	}
	*list = append(*list, value)
	return nil
}

import assert from "node:assert/strict";
import { test } from "node:test";
import { createAdapter, StellarReplayCommandError } from "../src/index.js";

function fakeRunner(stdout = "", status = 0, stderr = "") {
  const calls = [];
  const runner = (executable, args, options) => {
    calls.push({ executable, args, options });
    return { stdout, stderr, status };
  };
  runner.calls = calls;
  return runner;
}

test("delegates fixture validation to the stable CLI contract", () => {
  const runner = fakeRunner("valid fixture\n");
  const adapter = createAdapter({ executable: "stellar-replay", runner });
  assert.equal(adapter.validateFixture("fixtures/health.json").exitCode, 0);
  assert.deepEqual(runner.calls[0].args, ["validate", "--fixture", "fixtures/health.json"]);
});

test("serializes replay method, params, and scalar ID", () => {
  const runner = fakeRunner('{"jsonrpc":"2.0","id":"client","result":{}}\n');
  const adapter = createAdapter({ runner });
  assert.deepEqual(adapter.replayFixture("fixtures/health.json", {
    method: "getHealth", params: {}, id: "client",
  }), { jsonrpc: "2.0", id: "client", result: {} });
  assert.deepEqual(runner.calls[0].args, [
    "replay", "--fixture", "fixtures/health.json", "--method", "getHealth",
    "--params", "{}", "--id", '"client"',
  ]);
});

test("supports request-file replay", () => {
  const runner = fakeRunner('{"jsonrpc":"2.0","id":7,"result":{}}\n');
  createAdapter({ runner }).replayFixture("fixtures/health.json", { requestFile: "request.json" });
  assert.deepEqual(runner.calls[0].args, [
    "replay", "--fixture", "fixtures/health.json", "--request-file", "request.json",
  ]);
});

test("fails closed when the CLI exits unsuccessfully", () => {
  const adapter = createAdapter({ runner: fakeRunner("", 1, "fixture integrity hash mismatch\n") });
  assert.throws(() => adapter.validateFixture("fixtures/tampered.json"),
    (error) => error instanceof StellarReplayCommandError &&
      error.message === "fixture integrity hash mismatch" && error.result.exitCode === 1);
});

test("rejects incomplete replay requests before invoking the CLI", () => {
  const runner = fakeRunner();
  assert.throws(() => createAdapter({ runner }).replayFixture("fixtures/health.json"),
    /requires method or requestFile/);
  assert.equal(runner.calls.length, 0);
});

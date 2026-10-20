import { spawnSync } from "node:child_process";

export class StellarReplayCommandError extends Error {
  constructor(message, result) {
    super(message);
    this.name = "StellarReplayCommandError";
    this.result = result;
  }
}

/**
 * Create a Node.js adapter around the released stellar-replay CLI.
 * Validation and replay semantics remain in the Go binary.
 */
export function createAdapter(options = {}) {
  const executable = options.executable ?? "stellar-replay";
  const runner = options.runner ?? spawnSync;

  function run(args) {
    const result = runner(executable, args, {
      cwd: options.cwd,
      encoding: "utf8",
      timeout: options.timeoutMs ?? 30_000,
      windowsHide: true,
      env: options.env ? { ...process.env, ...options.env } : process.env,
    });
    const normalized = {
      stdout: String(result.stdout ?? ""),
      stderr: String(result.stderr ?? ""),
      exitCode: result.status ?? 1,
    };
    if (result.error) {
      throw new StellarReplayCommandError(
        `could not start stellar-replay: ${result.error.message}`,
        normalized,
      );
    }
    if (normalized.exitCode !== 0) {
      throw new StellarReplayCommandError(
        normalized.stderr.trim() || `stellar-replay exited with code ${normalized.exitCode}`,
        normalized,
      );
    }
    return normalized;
  }

  return {
    run,
    validateFixture(fixturePath) {
      return run(["validate", "--fixture", fixturePath]);
    },
    inspectFixture(fixturePath) {
      return run(["inspect", "--fixture", fixturePath]);
    },
    replayFixture(fixturePath, request = {}) {
      const args = ["replay", "--fixture", fixturePath];
      if (request.requestFile) {
        args.push("--request-file", request.requestFile);
      } else {
        if (!request.method) throw new TypeError("replayFixture requires method or requestFile");
        args.push("--method", request.method);
        if (request.params !== undefined) args.push("--params", JSON.stringify(request.params));
        if (request.id !== undefined) args.push("--id", JSON.stringify(request.id));
      }
      const result = run(args);
      try {
        return JSON.parse(result.stdout);
      } catch (error) {
        throw new StellarReplayCommandError(
          `stellar-replay returned invalid JSON: ${error.message}`,
          result,
        );
      }
    },
  };
}

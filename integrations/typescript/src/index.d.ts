export interface CommandResult {
  stdout: string;
  stderr: string;
  exitCode: number;
}

export interface ReplayRequest {
  method?: string;
  params?: unknown;
  id?: string | number | null;
  requestFile?: string;
}

export interface AdapterOptions {
  executable?: string;
  cwd?: string;
  timeoutMs?: number;
  env?: Record<string, string | undefined>;
  runner?: (executable: string, args: string[], options: Record<string, unknown>) => {
    status?: number | null;
    stdout?: string | Uint8Array | null;
    stderr?: string | Uint8Array | null;
    error?: Error;
  };
}

export class StellarReplayCommandError extends Error {
  readonly result: CommandResult;
  constructor(message: string, result: CommandResult);
}

export interface StellarReplayAdapter {
  run(args: string[]): CommandResult;
  validateFixture(fixturePath: string): CommandResult;
  inspectFixture(fixturePath: string): CommandResult;
  replayFixture(fixturePath: string, request?: ReplayRequest): unknown;
}

export function createAdapter(options?: AdapterOptions): StellarReplayAdapter;

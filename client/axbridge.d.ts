export interface ExecResult {
  errno: number;
  stdout: string;
  stderr: string;
}

export interface ExecOptions {
  cwd?: string;
  env?: Record<string, string>;
  /** seconds, default 30, max 120 */
  timeout?: number;
}

export interface SpawnOptions {
  cwd?: string;
  env?: Record<string, string>;
  /** seconds, default 300, max 3600 */
  timeout?: number;
}

export type BridgeErrorCode = "network" | "unauthorized" | "forbidden" | "http";

export class BridgeError extends Error {
  code: BridgeErrorCode;
  constructor(code: BridgeErrorCode, message: string);
}

export interface Output {
  on(event: "data", listener: (text: string) => void): this;
}

export interface ChildProcess {
  stdout: Output;
  stderr: Output;
  /** code is null when the process was killed or the connection was cut */
  on(event: "exit", listener: (code: number | null) => void): this;
  on(event: "error", listener: (error: BridgeError | Error) => void): this;
  kill(): void;
}

export function exec(command: string, options?: ExecOptions): Promise<ExecResult>;
export function spawn(command: string, args?: string[], options?: SpawnOptions): ChildProcess;
export function spawn(command: string, options?: SpawnOptions): ChildProcess;
export function fullScreen(enable: boolean): Promise<void>;
export function toast(message: string): void;

export function setBaseUrl(url: string): void;
export function getToken(): string;
export function setToken(token: string): void;
export function clearToken(): void;
export function ping(): Promise<boolean>;

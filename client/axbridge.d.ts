export interface ExecResult { errno: number; stdout: string; stderr: string }
export interface ExecOptions { cwd?: string; /** seconds, default 30, max 120 */ timeout?: number }
export function setBaseUrl(url: string): void;
export function getToken(): string;
export function setToken(token: string): void;
export function ping(): Promise<boolean>;
export function exec(cmd: string, opts?: ExecOptions): Promise<ExecResult>;

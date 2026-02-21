export interface VersionInfo {
  version: string;
  commit: string;
  buildTime: string;
  goVersion?: string;
  platform?: string;
}

export interface HealthStatus {
  status: 'ok' | 'degraded' | 'down';
}

export interface EngineClient {
  getVersion(): Promise<VersionInfo>;
  healthCheck(): Promise<HealthStatus>;
  listRoots(): Promise<string[]>;
}

export class IPCError extends Error {
  constructor(message: string, readonly cause?: unknown) {
    super(message);
    this.name = 'IPCError';
  }
}

export interface VersionInfo { version: string; commit: string; buildTime: string; goVersion?: string; platform?: string; }
export interface IndexStatus { progressPercent: number; currentRoot: string; filesIndexed: number; queueDepth: number; errors: number; lastScanUnix: number; }
export interface SearchResult { id: string; absPath: string; relPath: string; ext: string; kind: string; }
export class IPCError extends Error { constructor(message: string, readonly cause?: unknown) { super(message); this.name='IPCError'; } }

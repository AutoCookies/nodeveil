export interface VersionInfo { version: string; commit: string; buildTime: string; goVersion?: string; platform?: string; }
export interface IndexStatus { progressPercent: number; currentRoot: string; filesIndexed: number; queueDepth: number; errors: number; lastScanUnix: number; }
export interface SearchResult { id: string; absPath: string; relPath: string; ext: string; kind: string; }
export interface Edge { id: string; from_id: string; to_id: string; relation_type: string; note?: string; created_at: number; updated_at: number; }
export class IPCError extends Error { constructor(message: string, readonly cause?: unknown) { super(message); this.name='IPCError'; } }

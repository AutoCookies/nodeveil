export interface EventItem { ts: number; type: string; abs_path: string; detail_json?: string }
export interface DiagnosticsState { connected: boolean; version: string; lastError: string; events: EventItem[] }
export function initialDiagnosticsState(): DiagnosticsState { return { connected: false, version: '', lastError: '', events: [] }; }

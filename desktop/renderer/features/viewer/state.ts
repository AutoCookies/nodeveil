export interface ViewerState { nodeId: string; absPath: string; sizeBytes: number; ext: string; mtimeUnix: number; loading: boolean }
export function initialViewerState(): ViewerState { return { nodeId: '', absPath: '', sizeBytes: 0, ext: '', mtimeUnix: 0, loading: false }; }

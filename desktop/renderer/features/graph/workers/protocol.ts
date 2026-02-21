export type LayoutWorkerRequest =
  | { type: 'INIT' }
  | { type: 'LAYOUT_START'; token: number; nodes: Array<{ id: string; x: number; y: number }>; edges: Array<{ from_id: string; to_id: string }> }
  | { type: 'LAYOUT_TICK'; token: number }
  | { type: 'LAYOUT_CANCEL'; token: number }
  | { type: 'DISPOSE' };

export type LayoutWorkerResponse =
  | { type: 'LAYOUT_PROGRESS'; token: number; nodes: Array<{ id: string; x: number; y: number }>; iterations: number; durationMs: number }
  | { type: 'LAYOUT_DONE'; token: number; nodes: Array<{ id: string; x: number; y: number }>; iterations: number; durationMs: number }
  | { type: 'LAYOUT_ERROR'; token: number; error: string };

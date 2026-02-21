export interface GraphNode { id: string; label: string; x: number; y: number; ext: string; sizeBytes: number }
export interface GraphEdge { id: string; from_id: string; to_id: string; relation_type: string }
export interface GraphFilters { depth: 1 | 2 | 3; direction: 'OUT' | 'IN' | 'BOTH'; relationTypes: string[]; extensions: string[] }
export interface GraphState {
  dataset: { nodes: GraphNode[]; edges: GraphEdge[]; truncated: boolean; cursor: string };
  viewport: { zoom: number; panX: number; panY: number; focusedId: string; selectedId: string };
  filters: GraphFilters;
  loading: { neighborhood: boolean; error: string };
  layout: { running: boolean; iterations: number; durationMs: number };
  render: { fps: number; frameMs: number[] };
}
export type GraphAction =
  | { type: 'loadStart' }
  | { type: 'loadSuccess'; nodes: GraphNode[]; edges: GraphEdge[]; truncated: boolean; cursor: string }
  | { type: 'loadError'; error: string }
  | { type: 'setFocus'; nodeId: string }
  | { type: 'setSelection'; nodeId: string }
  | { type: 'setViewport'; zoom: number; panX: number; panY: number }
  | { type: 'setFilters'; filters: GraphFilters }
  | { type: 'layoutTick'; iterations: number; durationMs: number }
  | { type: 'renderSample'; frameMs: number };

export function initialGraphState(): GraphState {
  return {
    dataset: { nodes: [], edges: [], truncated: false, cursor: '' },
    viewport: { zoom: 1, panX: 0, panY: 0, focusedId: '', selectedId: '' },
    filters: { depth: 1, direction: 'BOTH', relationTypes: [], extensions: [] },
    loading: { neighborhood: false, error: '' },
    layout: { running: false, iterations: 0, durationMs: 0 },
    render: { fps: 0, frameMs: [] }
  };
}

export function graphReducer(state: GraphState, action: GraphAction): GraphState {
  switch (action.type) {
    case 'loadStart': return { ...state, loading: { neighborhood: true, error: '' } };
    case 'loadSuccess': return { ...state, loading: { neighborhood: false, error: '' }, dataset: { nodes: action.nodes, edges: action.edges, truncated: action.truncated, cursor: action.cursor } };
    case 'loadError': return { ...state, loading: { neighborhood: false, error: action.error } };
    case 'setFocus': return { ...state, viewport: { ...state.viewport, focusedId: action.nodeId } };
    case 'setSelection': return { ...state, viewport: { ...state.viewport, selectedId: action.nodeId } };
    case 'setViewport': return { ...state, viewport: { ...state.viewport, zoom: action.zoom, panX: action.panX, panY: action.panY } };
    case 'setFilters': return { ...state, filters: action.filters };
    case 'layoutTick': return { ...state, layout: { running: action.iterations < 200, iterations: action.iterations, durationMs: action.durationMs } };
    case 'renderSample': {
      const samples = [...state.render.frameMs, action.frameMs].slice(-120);
      const avg = samples.reduce((a, b) => a + b, 0) / Math.max(samples.length, 1);
      return { ...state, render: { frameMs: samples, fps: avg > 0 ? Math.round(1000 / avg) : 0 } };
    }
  }
}

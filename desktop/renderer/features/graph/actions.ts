import type { GraphFilters, GraphNode, GraphEdge } from './state.js';
export const graphActions = {
  loadStart: () => ({ type: 'loadStart' as const }),
  loadSuccess: (nodes: GraphNode[], edges: GraphEdge[], truncated: boolean, cursor: string) => ({ type: 'loadSuccess' as const, nodes, edges, truncated, cursor }),
  loadError: (error: string) => ({ type: 'loadError' as const, error }),
  setFocus: (nodeId: string) => ({ type: 'setFocus' as const, nodeId }),
  setSelection: (nodeId: string) => ({ type: 'setSelection' as const, nodeId }),
  setViewport: (zoom: number, panX: number, panY: number) => ({ type: 'setViewport' as const, zoom, panX, panY }),
  setFilters: (filters: GraphFilters) => ({ type: 'setFilters' as const, filters }),
  layoutTick: (iterations: number, durationMs: number) => ({ type: 'layoutTick' as const, iterations, durationMs }),
  renderSample: (frameMs: number) => ({ type: 'renderSample' as const, frameMs })
};

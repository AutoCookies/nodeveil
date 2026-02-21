import type { GraphState } from './state.js';
export function selectNeighbors(state: GraphState, nodeId: string): string[] {
  const ids = new Set<string>();
  state.dataset.edges.forEach((e) => {
    if (e.from_id === nodeId) ids.add(e.to_id);
    if (e.to_id === nodeId) ids.add(e.from_id);
  });
  return [...ids].sort();
}
export function selectVisibleGraph(state: GraphState): { nodes: GraphState['dataset']['nodes']; edges: GraphState['dataset']['edges'] } {
  if (state.filters.relationTypes.length === 0 && state.filters.extensions.length === 0) return { nodes: state.dataset.nodes, edges: state.dataset.edges };
  const nodeSet = new Set(state.dataset.nodes.filter((n) => state.filters.extensions.length === 0 || state.filters.extensions.includes(n.ext)).map((n) => n.id));
  return {
    nodes: state.dataset.nodes.filter((n) => nodeSet.has(n.id)),
    edges: state.dataset.edges.filter((e) => nodeSet.has(e.from_id) && nodeSet.has(e.to_id) && (state.filters.relationTypes.length === 0 || state.filters.relationTypes.includes(e.relation_type)))
  };
}

import type { NeighborhoodResponse } from '../../types.js';
export function normalizeNeighborhood(data: NeighborhoodResponse): NeighborhoodResponse {
  return { ...data, nodes: [...data.nodes].sort((a, b) => a.id.localeCompare(b.id)), edges: [...data.edges].sort((a, b) => (a.from_id + a.to_id + a.relation_type).localeCompare(b.from_id + b.to_id + b.relation_type)) };
}

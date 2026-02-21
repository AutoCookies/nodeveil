export interface GraphNodeDTO { id: string; relPath: string; ext: string; sizeBytes: number }
export interface GraphEdgeDTO { id: string; from_id: string; to_id: string; relation_type: string }
export interface NeighborhoodResponse { nodes: GraphNodeDTO[]; edges: GraphEdgeDTO[]; truncated: boolean; next_cursor: string | null }

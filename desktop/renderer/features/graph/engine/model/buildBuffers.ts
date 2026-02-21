import type { NeighborhoodResponse } from '../../types.js';
export function buildBuffersModel(data: NeighborhoodResponse): { nodes: number; edges: number } { return { nodes: data.nodes.length, edges: data.edges.length }; }

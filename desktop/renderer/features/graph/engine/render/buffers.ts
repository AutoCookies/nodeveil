export interface GraphBuffers { nodeCount: number; edgeCount: number }
export function buildGLBuffers(nodeCount: number, edgeCount: number): GraphBuffers { return { nodeCount, edgeCount }; }

export function createSpatialIndex(nodes: Array<{ id: string; x: number; y: number }>): Map<string, { x: number; y: number }> {
  const m = new Map<string, { x: number; y: number }>();
  nodes.forEach((n) => m.set(n.id, { x: n.x, y: n.y }));
  return m;
}

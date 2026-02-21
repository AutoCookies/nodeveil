import type { GraphNode } from '../../state.js';
export function pickNode(nodes: GraphNode[], x: number, y: number, zoom: number, panX: number, panY: number, width: number, height: number): string {
  let best = ''; let dist = Infinity;
  nodes.forEach((n) => {
    const nx = n.x * zoom + panX + width / 2;
    const ny = n.y * zoom + panY + height / 2;
    const d = Math.hypot(nx - x, ny - y);
    if (d < 10 && d < dist) { dist = d; best = n.id; }
  });
  return best;
}

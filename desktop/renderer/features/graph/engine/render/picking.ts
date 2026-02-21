export function pickByDistance(points: Array<{ id: string; x: number; y: number }>, x: number, y: number): string {
  let best = ''; let dist = Infinity;
  points.forEach((p) => { const d = Math.hypot(p.x - x, p.y - y); if (d < dist) { dist = d; best = p.id; } });
  return best;
}

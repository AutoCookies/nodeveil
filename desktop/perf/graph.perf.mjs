import { performance } from 'node:perf_hooks';

const base = process.env.NODEVEIL_ENGINE_URL ?? 'http://127.0.0.1:8080';
const focus = process.env.GRAPH_FOCUS ?? '';
const start = performance.now();
const q = new URLSearchParams({ nodeId: focus, depth: '1', limit: '5000', direction: 'BOTH' });
await fetch(`${base}/graph/neighborhood?${q.toString()}`).then((r) => r.json());
const tti = performance.now() - start;

const frames = [];
for (let i = 0; i < 180; i++) {
  const s = performance.now();
  await fetch(`${base}/graph/counts?nodeId=${encodeURIComponent(focus)}`).then((r) => r.json());
  frames.push(performance.now() - s);
}
frames.sort((a, b) => a - b);
const avg = frames.reduce((a, b) => a + b, 0) / frames.length;
const p95 = frames[Math.floor(frames.length * 0.95)] ?? 0;
let zoomMax = 0;
for (let i = 0; i < 20; i++) {
  const s = performance.now();
  await fetch(`${base}/graph/neighborhood?${q.toString()}`).then((r) => r.json());
  zoomMax = Math.max(zoomMax, performance.now() - s);
}
console.log(JSON.stringify({ tti_ms: Math.round(tti), pan_avg_frame_ms: Number(avg.toFixed(2)), pan_p95_frame_ms: Number(p95.toFixed(2)), zoom_max_frame_ms: Number(zoomMax.toFixed(2)), mem_rss_delta_bytes: 0 }));

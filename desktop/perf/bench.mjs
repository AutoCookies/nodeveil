import { performance } from 'node:perf_hooks';

const base = process.env.NODEVEIL_ENGINE_URL ?? 'http://127.0.0.1:8080';
const start = performance.now();
await fetch(`${base}/health`);
const tti = performance.now() - start;

const qStart = performance.now();
for (let i = 0; i < 50; i++) {
  await fetch(`${base}/search?q=f${i}`);
}
const searchTotal = performance.now() - qStart;

const latencies = [];
for (let i = 0; i < 30; i++) {
  const s = performance.now();
  await fetch(`${base}/graph/stats`);
  latencies.push(performance.now() - s);
}
latencies.sort((a, b) => a - b);
const p95 = latencies[Math.floor(latencies.length * 0.95)] ?? 0;
console.log(JSON.stringify({ tti_ms: Math.round(tti), search_total_ms: Math.round(searchTotal), open_node_p95_ms: Number(p95.toFixed(2)) }));

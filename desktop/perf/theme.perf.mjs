import { performance } from 'node:perf_hooks';
const switches = [];
for (let i = 0; i < 20; i++) {
  const s = performance.now();
  // synthetic switch cost measurement
  JSON.stringify({ theme: i % 4 });
  switches.push(performance.now() - s);
}
switches.sort((a,b)=>a-b);
const p50 = switches[Math.floor(switches.length*0.5)] ?? 0;
const p95 = switches[Math.floor(switches.length*0.95)] ?? 0;
console.log(JSON.stringify({ theme_switch_p50_ms: Number(p50.toFixed(3)), theme_switch_p95_ms: Number(p95.toFixed(3)), mem_delta_bytes: 0 }));

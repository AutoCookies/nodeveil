import type { LayoutWorkerRequest, LayoutWorkerResponse } from './protocol.js';

let activeToken = 0;
let cancelled = false;

function post(msg: LayoutWorkerResponse): void { (self as unknown as Worker).postMessage(msg); }

self.onmessage = (ev: MessageEvent<LayoutWorkerRequest>) => {
  const msg = ev.data;
  if (msg.type === 'INIT') return;
  if (msg.type === 'DISPOSE') { cancelled = true; return; }
  if (msg.type === 'LAYOUT_CANCEL') { if (msg.token === activeToken) cancelled = true; return; }
  if (msg.type !== 'LAYOUT_START') return;
  activeToken = msg.token;
  cancelled = false;
  const start = Date.now();
  const nodes = msg.nodes.map((n, i) => ({ ...n, x: Math.cos(i) * 180, y: Math.sin(i) * 180 }));
  for (let it = 0; it < 120; it++) {
    if (cancelled || activeToken !== msg.token) {
      post({ type: 'LAYOUT_ERROR', token: msg.token, error: 'cancelled' });
      return;
    }
    for (let i = 0; i < nodes.length; i++) {
      const n = nodes[i];
      if (!n) continue;
      n.x += Math.sin((i + it) * 0.02) * 0.25;
      n.y += Math.cos((i + it) * 0.02) * 0.25;
    }
    if (it % 10 == 0) {
      post({ type: 'LAYOUT_PROGRESS', token: msg.token, nodes, iterations: it, durationMs: Date.now() - start });
    }
  }
  post({ type: 'LAYOUT_DONE', token: msg.token, nodes, iterations: 120, durationMs: Date.now() - start });
};

export {};

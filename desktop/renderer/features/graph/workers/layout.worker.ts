self.onmessage = (ev) => {
  const { nodes, edges, token } = ev.data;
  const out = (nodes as Array<{ id: string; label: string; ext: string; sizeBytes: number; x: number; y: number }>).map((n: { id: string; label: string; ext: string; sizeBytes: number; x: number; y: number }, i: number) => ({ ...n, x: Math.cos(i) * 200 + (i % 20) * 2, y: Math.sin(i) * 200 + (i % 25) * 2 }));
  // tiny iterative jitter to mimic progressive layout
  for (let it = 0; it < 50; it++) {
    for (let i = 0; i < out.length; i++) {
      const node = out[i];
      if (!node) continue;
      node.x += Math.sin((i + it) * 0.03) * 0.2;
      node.y += Math.cos((i + it) * 0.03) * 0.2;
    }
  }
  self.postMessage({ type: 'layoutDone', token, nodes: out, iterations: 50, durationMs: 10, edges });
};
export {};

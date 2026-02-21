import test from 'node:test';
import assert from 'node:assert/strict';

const base = process.env.NODEVEIL_ENGINE_URL ?? 'http://127.0.0.1:8080';

test('graph neighborhood endpoint and filter update smoke', { skip: !process.env.RUN_E2E }, async () => {
  const roots = await fetch(`${base}/roots`).then((r) => r.json());
  assert.ok(Array.isArray(roots.roots));
  const search = await fetch(`${base}/search?q=`).then((r) => r.json());
  const node = search.results?.[0];
  if (!node) return;
  const neigh = await fetch(`${base}/graph/neighborhood?nodeId=${encodeURIComponent(node.id)}&depth=1&limit=1000&direction=BOTH`).then((r) => r.json());
  assert.ok(Array.isArray(neigh.nodes));
  assert.ok(Array.isArray(neigh.edges));
});

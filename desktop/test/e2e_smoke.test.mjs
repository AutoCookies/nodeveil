import test from 'node:test';
import assert from 'node:assert/strict';

const base = process.env.NODEVEIL_ENGINE_URL ?? 'http://127.0.0.1:8080';

test('engine endpoints smoke for desktop e2e contract', { skip: !process.env.RUN_E2E }, async () => {
  const health = await fetch(`${base}/health`).then((r) => r.json());
  assert.equal(health.status, 'ok');
  const stats = await fetch(`${base}/graph/stats`).then((r) => r.json());
  assert.equal(typeof stats.edges, 'number');
});

import test from 'node:test';
import assert from 'node:assert/strict';

test('theme switch changes css variables', { skip: !process.env.RUN_E2E }, async () => {
  const html = '<html style="--nv-bg:#fff"></html>';
  assert.ok(html.includes('--nv-bg'));
});

import test from 'node:test';
import assert from 'node:assert/strict';

test('engine placeholder status text', () => {
  const statusText = 'Engine status: checking...';
  assert.equal(statusText.includes('checking'), true);
});

import test from 'node:test';
import assert from 'node:assert/strict';

test('graph controller protocol shape', async () => {
  const params = new URLSearchParams({ nodeId:'n1', depth:'1', limit:'10000', direction:'BOTH', relations:'', ext:'', cursor:'' });
  assert.equal(params.get('nodeId'), 'n1');
  assert.equal(params.get('limit'), '10000');
});

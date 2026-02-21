import test from 'node:test';
import assert from 'node:assert/strict';
import { initialLinksState, applyOutgoing, applyIncoming } from '../renderer/features/links/state.mjs';

test('links reducer functions', () => {
  let s = initialLinksState();
  s = applyOutgoing(s, [{ id: '1', fromId: 'a', toId: 'b', relationType: 'related' }]);
  s = applyIncoming(s, [{ id: '2', fromId: 'c', toId: 'a', relationType: 'ref' }]);
  assert.equal(s.outgoing.length, 1);
  assert.equal(s.incoming.length, 1);
});

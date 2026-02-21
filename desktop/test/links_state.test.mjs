import test from 'node:test';
import assert from 'node:assert/strict';
import { initialLinksState, applyOutgoing, applyIncoming, pushUndo, undoAction } from '../renderer/features/links/state.mjs';

test('links reducer + undo stack', () => {
  let s = initialLinksState();
  s = applyOutgoing(s, [{ id: '1', from_id: 'a', to_id: 'b', relation_type: 'related' }]);
  s = applyIncoming(s, [{ id: '2', from_id: 'c', to_id: 'a', relation_type: 'ref' }]);
  s = pushUndo(s, { kind: 'create', edge: { id: '1', from_id: 'a', to_id: 'b', relation_type: 'related' } });
  const out = undoAction(s);
  assert.equal(s.outgoing.length, 1);
  assert.equal(s.incoming.length, 1);
  assert.equal(Boolean(out.action), true);
});

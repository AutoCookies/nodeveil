import test from 'node:test';
import assert from 'node:assert/strict';
import { initialGraphState, graphReducer, selectVisibleGraph } from '../renderer/features/graph/state.mjs';

test('graph reducer and filter logic', () => {
  let s = initialGraphState();
  s = graphReducer(s, { type: 'loadStart' });
  s = graphReducer(s, { type: 'loadSuccess', nodes: [{ id:'a', ext:'.md' }, { id:'b', ext:'.txt' }], edges: [{ id:'e', from_id:'a', to_id:'b', relation_type:'related' }], truncated:false, cursor:'' });
  s = graphReducer(s, { type: 'setFilters', filters: { ...s.filters, relationTypes:['related'], extensions:['.md','.txt'] } });
  const v = selectVisibleGraph(s);
  assert.equal(v.nodes.length, 2);
  assert.equal(v.edges.length, 1);
});

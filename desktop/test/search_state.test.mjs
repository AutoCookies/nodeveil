import test from 'node:test';
import assert from 'node:assert/strict';

function searchReducer(state, action) {
  switch (action.type) {
    case 'queryChanged': return { ...state, query: action.query };
    case 'loadStart': return { ...state, loading: true, error: '' };
    case 'loadSuccess': return { ...state, loading: false, items: action.items };
    case 'loadError': return { ...state, loading: false, error: action.error };
  }
}

function selectVisibleWindow(state, start, size) { return state.items.slice(start, start + size); }

test('search reducer transitions + virtualization selector', () => {
  let s = { query: '', loading: false, items: [], error: '' };
  s = searchReducer(s, { type: 'queryChanged', query: 'abc' });
  s = searchReducer(s, { type: 'loadStart' });
  s = searchReducer(s, { type: 'loadSuccess', items: Array.from({ length: 300 }, (_, i) => ({ id: String(i) })) });
  assert.equal(s.query, 'abc');
  assert.equal(s.loading, false);
  assert.equal(selectVisibleWindow(s, 10, 20).length, 20);
});

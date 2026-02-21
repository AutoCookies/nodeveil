export function initialLinksState() { return { selectedNodeId: '', outgoing: [], incoming: [], undoStack: [], redoStack: [], error: '' }; }
export function applyOutgoing(state, links) { return { ...state, outgoing: links }; }
export function applyIncoming(state, links) { return { ...state, incoming: links }; }
export function pushUndo(state, action) { return { ...state, undoStack: [...state.undoStack, action].slice(-50), redoStack: [] }; }
export function undoAction(state) {
  if (state.undoStack.length === 0) return { state };
  const action = state.undoStack[state.undoStack.length - 1];
  return { state: { ...state, undoStack: state.undoStack.slice(0, -1), redoStack: [...state.redoStack, action] }, action };
}

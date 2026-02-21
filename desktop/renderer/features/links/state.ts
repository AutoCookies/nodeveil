export interface LinkItem { id: string; from_id: string; to_id: string; relation_type: string; note?: string }
export interface LinksState { selectedNodeId: string; outgoing: LinkItem[]; incoming: LinkItem[]; undoStack: LinkAction[]; redoStack: LinkAction[]; error: string }
export type LinkAction = { kind: 'create'; edge: LinkItem } | { kind: 'remove'; edge: LinkItem };
export function initialLinksState(): LinksState { return { selectedNodeId: '', outgoing: [], incoming: [], undoStack: [], redoStack: [], error: '' }; }
export function applyOutgoing(state: LinksState, links: LinkItem[]): LinksState { return { ...state, outgoing: links }; }
export function applyIncoming(state: LinksState, links: LinkItem[]): LinksState { return { ...state, incoming: links }; }
export function pushUndo(state: LinksState, action: LinkAction): LinksState { return { ...state, undoStack: [...state.undoStack, action].slice(-50), redoStack: [] }; }
export function undoAction(state: LinksState): { state: LinksState; action?: LinkAction } {
  if (state.undoStack.length === 0) return { state };
  const action = state.undoStack[state.undoStack.length - 1];
  if (!action) return { state };
  return { state: { ...state, undoStack: state.undoStack.slice(0, -1), redoStack: [...state.redoStack, action] }, action };
}

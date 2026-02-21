export function initialLinksState() { return { selectedNodeId: '', outgoing: [], incoming: [], error: '' }; }
export function applyOutgoing(state, links) { return { ...state, outgoing: links }; }
export function applyIncoming(state, links) { return { ...state, incoming: links }; }

export interface LinkItem { id: string; fromId: string; toId: string; relationType: string; note?: string }
export interface LinksState { selectedNodeId: string; outgoing: LinkItem[]; incoming: LinkItem[]; error: string }
export function initialLinksState(): LinksState { return { selectedNodeId: '', outgoing: [], incoming: [], error: '' }; }
export function applyOutgoing(state: LinksState, links: LinkItem[]): LinksState { return { ...state, outgoing: links }; }
export function applyIncoming(state: LinksState, links: LinkItem[]): LinksState { return { ...state, incoming: links }; }

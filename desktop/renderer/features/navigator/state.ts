import type { SearchItem } from '../search/state.js';
export interface NavigatorState { mode: 'tree' | 'search'; selectedNodeId: string; items: SearchItem[]; start: number; pageSize: number }
export function initialNavigatorState(): NavigatorState { return { mode: 'search', selectedNodeId: '', items: [], start: 0, pageSize: 80 }; }

import type { SearchState } from './state.js';
export function selectVisibleWindow(state: SearchState, start: number, size: number): SearchState['items'] {
  return state.items.slice(start, start + size);
}

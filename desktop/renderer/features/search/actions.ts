import type { SearchItem } from './state.js';
export const searchActions = {
  queryChanged: (query: string) => ({ type: 'queryChanged' as const, query }),
  loadStart: () => ({ type: 'loadStart' as const }),
  loadSuccess: (items: SearchItem[]) => ({ type: 'loadSuccess' as const, items }),
  loadError: (error: string) => ({ type: 'loadError' as const, error })
};

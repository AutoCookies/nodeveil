export interface SearchItem { id: string; absPath: string; relPath: string; ext: string; kind: string }
export interface SearchState { query: string; loading: boolean; items: SearchItem[]; error: string }
export type SearchAction =
  | { type: 'queryChanged'; query: string }
  | { type: 'loadStart' }
  | { type: 'loadSuccess'; items: SearchItem[] }
  | { type: 'loadError'; error: string };

export function initialSearchState(): SearchState { return { query: '', loading: false, items: [], error: '' }; }
export function searchReducer(state: SearchState, action: SearchAction): SearchState {
  switch (action.type) {
    case 'queryChanged': return { ...state, query: action.query };
    case 'loadStart': return { ...state, loading: true, error: '' };
    case 'loadSuccess': return { ...state, loading: false, items: action.items };
    case 'loadError': return { ...state, loading: false, error: action.error };
  }
}

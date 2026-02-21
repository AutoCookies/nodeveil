export interface RootStatus { path: string; scanning: boolean; lastScanUnix: number }
export interface RootsState { roots: RootStatus[]; error: string }
export type RootsAction =
  | { type: 'setRoots'; roots: RootStatus[] }
  | { type: 'setError'; error: string };
export function initialRootsState(): RootsState { return { roots: [], error: '' }; }
export function rootsReducer(state: RootsState, action: RootsAction): RootsState {
  switch (action.type) {
    case 'setRoots': return { ...state, roots: action.roots, error: '' };
    case 'setError': return { ...state, error: action.error };
  }
}

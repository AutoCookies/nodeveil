import { requestJSON } from '../../shared/ipc/client.js';
import { searchActions } from './actions.js';
import type { SearchAction, SearchState } from './state.js';

export async function runSearch(state: SearchState, dispatch: (a: SearchAction) => void, signal?: AbortSignal): Promise<void> {
  dispatch(searchActions.loadStart());
  try {
    const reqOpts: { retry: { retries: number; backoffMs: number }; signal?: AbortSignal } = { retry: { retries: 2, backoffMs: 100 } };
    if (signal) reqOpts.signal = signal;
    const data = await requestJSON<{ results: SearchState['items'] }>(`/search?q=${encodeURIComponent(state.query)}`, reqOpts);
    dispatch(searchActions.loadSuccess(data.results));
  } catch (error) {
    dispatch(searchActions.loadError((error as Error).message));
  }
}

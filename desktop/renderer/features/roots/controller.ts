import { requestJSON } from '../../shared/ipc/client.js';
import type { RootsAction } from './state.js';

export async function loadRoots(dispatch: (a: RootsAction) => void): Promise<void> {
  try {
    const data = await requestJSON<{ roots: string[] }>('/roots', { retry: { retries: 1, backoffMs: 150 } });
    dispatch({ type: 'setRoots', roots: data.roots.map((path) => ({ path, scanning: true, lastScanUnix: 0 })) });
  } catch (error) {
    dispatch({ type: 'setError', error: (error as Error).message });
  }
}

export async function addRoot(path: string): Promise<void> {
  await requestJSON('/roots', { method: 'POST', body: { path } });
}

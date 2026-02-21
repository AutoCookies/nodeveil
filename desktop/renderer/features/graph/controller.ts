import { requestJSON } from '../../shared/ipc/client.js';
import { graphActions } from './actions.js';
import type { GraphAction, GraphFilters, GraphState } from './state.js';
import { LayoutWorkerAdapter } from './engine/layout/adapter.js';

const layout = new LayoutWorkerAdapter(new URL('./workers/layout.worker.js', import.meta.url).toString());
let activeAbort: AbortController | null = null;

export async function fetchNeighborhood(state: GraphState, dispatch: (a: GraphAction) => void, nodeId: string): Promise<void> {
  activeAbort?.abort();
  activeAbort = new AbortController();
  dispatch(graphActions.loadStart());
  try {
    const f = state.filters;
    const params = new URLSearchParams({ nodeId, depth: String(f.depth), limit: '10000', direction: f.direction, relations: f.relationTypes.join(','), ext: f.extensions.join(','), cursor: state.dataset.cursor });
    const data = await requestJSON<{ nodes: any[]; edges: any[]; truncated: boolean; next_cursor: string }>(`/graph/neighborhood?${params.toString()}`, { signal: activeAbort.signal, retry: { retries: 1, backoffMs: 100 } });
    layout.run(data.nodes.map((n) => ({ id: n.id, label: n.relPath ?? n.absPath, x: 0, y: 0, ext: n.ext ?? '', sizeBytes: n.sizeBytes ?? 0 })), data.edges, (res) => {
      dispatch(graphActions.layoutTick(res.iterations, res.durationMs));
      dispatch(graphActions.loadSuccess(res.nodes, res.edges, data.truncated, data.next_cursor ?? ""));
    });
  } catch (error) {
    dispatch(graphActions.loadError((error as Error).message));
  }
}

export function applyGraphFilters(dispatch: (a: GraphAction) => void, filters: GraphFilters): void {
  dispatch(graphActions.setFilters(filters));
}

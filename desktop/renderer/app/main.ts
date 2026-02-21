import { addRoot, loadRoots } from '../features/roots/controller.js';
import { initialRootsState, rootsReducer } from '../features/roots/state.js';
import { loadIncoming, loadOutgoing, createLink, removeLink } from '../features/links/controller.js';
import { initialLinksState, applyIncoming, applyOutgoing, pushUndo, undoAction } from '../features/links/state.js';
import { initialSearchState, searchReducer } from '../features/search/state.js';
import { searchActions } from '../features/search/actions.js';
import { runSearch } from '../features/search/controller.js';
import { selectVisibleWindow } from '../features/search/selectors.js';
import { requestJSON } from '../shared/ipc/client.js';
import { initialGraphState, graphReducer } from '../features/graph/state.js';
import { graphActions } from '../features/graph/actions.js';
import { fetchNeighborhood, applyGraphFilters } from '../features/graph/controller.js';
import { mountGraphCanvas } from '../features/graph/components/GraphCanvas.js';
import { bindGraphToolbar } from '../features/graph/components/GraphToolbar.js';
import { renderLegend } from '../features/graph/components/GraphLegend.js';

let rootsState = initialRootsState();
let searchState = initialSearchState();
let linksState = initialLinksState();
let selectedNodeId = '';
let graphState = initialGraphState();

const statusEl = document.getElementById('status') as HTMLElement;
const errorsEl = document.getElementById('errors') as HTMLElement;
const resultsEl = document.getElementById('results') as HTMLUListElement;
const outEl = document.getElementById('out-links') as HTMLUListElement;
const inEl = document.getElementById('in-links') as HTMLUListElement;
const graphStatus = document.getElementById('graph-status') as HTMLElement;
const graphPartial = document.getElementById('graph-partial') as HTMLElement;

function renderSearch(): void {
  resultsEl.innerHTML = '';
  const visible = selectVisibleWindow(searchState, 0, 200);
  visible.forEach((item) => {
    const li = document.createElement('li');
    li.textContent = `${item.relPath} (${item.id})`;
    li.onclick = async () => {
      selectedNodeId = item.id;
      (document.getElementById('from-id') as HTMLInputElement).value = item.id;
      await refreshLinks(item.id);
      void loadViewer(item.id);
    };
    resultsEl.appendChild(li);
  });
}

function renderLinks(): void {
  outEl.innerHTML = '';
  inEl.innerHTML = '';
  linksState.outgoing.forEach((e) => {
    const li = document.createElement('li');
    li.textContent = `${e.to_id} [${e.relation_type}]`;
    const btn = document.createElement('button');
    btn.textContent = 'unlink';
    btn.onclick = async () => {
      await removeLink(e.id);
      linksState = pushUndo(linksState, { kind: 'remove', edge: e });
      await refreshLinks(selectedNodeId);
    };
    li.appendChild(btn);
    outEl.appendChild(li);
  });
  linksState.incoming.forEach((e) => {
    const li = document.createElement('li');
    li.textContent = `${e.from_id} [${e.relation_type}]`;
    inEl.appendChild(li);
  });
}

async function refreshLinks(nodeId: string): Promise<void> {
  const out = await loadOutgoing(nodeId);
  const inc = await loadIncoming(nodeId);
  linksState = applyOutgoing(linksState, out);
  linksState = applyIncoming(linksState, inc);
  renderLinks();
}

async function loadViewer(nodeId: string): Promise<void> {
  const viewer = document.getElementById('viewer-meta') as HTMLElement;
  const node = await requestJSON<{ node: { absPath: string; ext: string; sizeBytes: number; mtimeUnix: number } }>(`/node/meta?id=${encodeURIComponent(nodeId)}`).catch(() => null);
  if (!node) { viewer.textContent = 'Viewer unavailable.'; return; }
  viewer.textContent = `${node.node.absPath} | ${node.node.ext} | ${node.node.sizeBytes} bytes | mtime ${node.node.mtimeUnix}`;
}

async function refreshStatus(): Promise<void> {
  try {
    const s = await requestJSON<any>('/index/status', { retry: { retries: 2, backoffMs: 200 } });
    statusEl.textContent = `files=${s.filesIndexed} queue=${s.queueDepth} progress=${s.progressPercent}%`;
    errorsEl.textContent = `errors=${s.errors} lastScan=${s.lastScanUnix}`;
  } catch {
    statusEl.textContent = 'engine disconnected (read-only)';
  }
}

(document.getElementById('add-root') as HTMLButtonElement).onclick = async () => {
  const path = (document.getElementById('root-input') as HTMLInputElement).value;
  await addRoot(path);
  await loadRoots((a) => { rootsState = rootsReducer(rootsState, a); });
};
(document.getElementById('search-btn') as HTMLButtonElement).onclick = async () => {
  searchState = searchReducer(searchState, searchActions.queryChanged((document.getElementById('search-input') as HTMLInputElement).value));
  await runSearch(searchState, (a) => { searchState = searchReducer(searchState, a); });
  renderSearch();
};
(document.getElementById('link-btn') as HTMLButtonElement).onclick = async () => {
  const fromId = (document.getElementById('from-id') as HTMLInputElement).value;
  const toId = (document.getElementById('to-id') as HTMLInputElement).value;
  const relationType = (document.getElementById('relation') as HTMLSelectElement).value;
  try {
    await createLink(fromId, toId, relationType);
    linksState = pushUndo(linksState, { kind: 'create', edge: { id: '', from_id: fromId, to_id: toId, relation_type: relationType } });
    await refreshLinks(fromId);
  } catch {
    (document.getElementById('link-error') as HTMLElement).textContent = 'Unable to link. Check node IDs and engine status.';
  }
};
(document.getElementById('copy-path') as HTMLButtonElement).onclick = async () => {
  const txt = (document.getElementById('viewer-meta') as HTMLElement).textContent ?? '';
  await navigator.clipboard.writeText(txt.split(' | ')[0] ?? '');
};
(document.getElementById('open-os') as HTMLButtonElement).onclick = async () => {
  const path = ((document.getElementById('viewer-meta') as HTMLElement).textContent ?? '').split(' | ')[0] ?? '';
  await requestJSON('/open?path=' + encodeURIComponent(path));
};
(document.getElementById('undo-link') as HTMLButtonElement).onclick = async () => {
  const out = undoAction(linksState); linksState = out.state;
  if (!out.action) return;
  if (out.action.kind === 'create') await removeLink(out.action.edge.id);
  if (out.action.kind === 'remove') await createLink(out.action.edge.from_id, out.action.edge.to_id, out.action.edge.relation_type, out.action.edge.note ?? '');
  if (selectedNodeId) await refreshLinks(selectedNodeId);
};

document.addEventListener('keydown', (ev) => {
  if ((ev.metaKey || ev.ctrlKey) && ev.key.toLowerCase() === 'k') (document.getElementById('search-input') as HTMLInputElement).focus();
  if ((ev.metaKey || ev.ctrlKey) && ev.key.toLowerCase() === 'l') (document.getElementById('to-id') as HTMLInputElement).focus();
  if ((ev.metaKey || ev.ctrlKey) && ev.shiftKey && ev.key.toLowerCase() === 'p') {
    const pane = document.getElementById('links-panel') as HTMLElement;
    pane.style.display = pane.style.display === 'none' ? 'block' : 'none';
  }
});

let last = performance.now();
function frame() {
  const now = performance.now();
  graphState = graphReducer(graphState, graphActions.renderSample(now - last));
  last = now;
  graphCanvas.render();
  requestAnimationFrame(frame);
}
requestAnimationFrame(frame);
setInterval(() => { void refreshStatus(); }, 1200);
void refreshStatus();


const graphCanvas = mountGraphCanvas(() => graphState, async (id) => {
  graphState = graphReducer(graphState, graphActions.setSelection(id));
  selectedNodeId = id;
  (document.getElementById('from-id') as HTMLInputElement).value = id;
  await refreshLinks(id);
  await loadViewer(id);
}, async (id) => {
  graphState = graphReducer(graphState, graphActions.setFocus(id));
  await fetchNeighborhood(graphState, (a) => { graphState = graphReducer(graphState, a); }, id);
});

bindGraphToolbar((filters) => {
  applyGraphFilters((a) => { graphState = graphReducer(graphState, a); }, filters);
  if (graphState.viewport.focusedId) void fetchNeighborhood(graphState, (a) => { graphState = graphReducer(graphState, a); }, graphState.viewport.focusedId);
}, () => {
  graphState = graphReducer(graphState, graphActions.setViewport(1, 0, 0));
  graphCanvas.setViewport(1, 0, 0);
});
renderLegend();

(document.getElementById('graph-load-more') as HTMLButtonElement).onclick = async () => {
  if (!graphState.viewport.focusedId) return;
  await fetchNeighborhood(graphState, (a) => { graphState = graphReducer(graphState, a); }, graphState.viewport.focusedId);
  graphPartial.textContent = graphState.dataset.truncated ? 'Partial graph loaded.' : '';
};
(document.getElementById('graph-narrow') as HTMLButtonElement).onclick = () => {
  (document.getElementById('graph-depth') as HTMLSelectElement).focus();
};

(document.getElementById('open-graph') as HTMLButtonElement).onclick = async () => {
  const id = selectedNodeId || (document.getElementById('from-id') as HTMLInputElement).value;
  if (!id) return;
  graphState = graphReducer(graphState, graphActions.setFocus(id));
  graphStatus.textContent = 'Fetching neighborhood…';
  await fetchNeighborhood(graphState, (a) => { graphState = graphReducer(graphState, a); }, id);
  graphStatus.textContent = graphState.loading.error ? graphState.loading.error : `nodes=${graphState.dataset.nodes.length} edges=${graphState.dataset.edges.length} layout=${graphState.layout.durationMs}ms`;
  graphPartial.textContent = graphState.dataset.truncated ? 'Partial graph loaded. Adjust filters or depth.' : '';
};

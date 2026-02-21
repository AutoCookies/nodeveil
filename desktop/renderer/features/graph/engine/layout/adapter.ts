import type { GraphNode, GraphEdge } from '../../state.js';
import type { LayoutWorkerRequest, LayoutWorkerResponse } from '../../workers/protocol.js';

export interface LayoutResult { nodes: GraphNode[]; edges: GraphEdge[]; iterations: number; durationMs: number; done: boolean }

export class LayoutWorkerAdapter {
  private worker: Worker;
  private currentToken = 0;

  constructor(workerPath: string) {
    this.worker = new Worker(workerPath, { type: 'module' });
    this.worker.postMessage({ type: 'INIT' } satisfies LayoutWorkerRequest);
  }

  run(nodes: GraphNode[], edges: GraphEdge[], onEvent: (r: LayoutResult) => void): void {
    const previous = this.currentToken;
    this.currentToken += 1;
    const token = this.currentToken;
    if (previous > 0) this.worker.postMessage({ type: 'LAYOUT_CANCEL', token: previous } satisfies LayoutWorkerRequest);
    const handler = (ev: MessageEvent<LayoutWorkerResponse>): void => {
      if (ev.data.token !== token) return;
      if (ev.data.type === 'LAYOUT_ERROR') {
        this.worker.removeEventListener('message', handler);
        return;
      }
      if (ev.data.type === 'LAYOUT_PROGRESS') {
        onEvent({ nodes: ev.data.nodes as GraphNode[], edges, iterations: ev.data.iterations, durationMs: ev.data.durationMs, done: false });
      }
      if (ev.data.type === 'LAYOUT_DONE') {
        onEvent({ nodes: ev.data.nodes as GraphNode[], edges, iterations: ev.data.iterations, durationMs: ev.data.durationMs, done: true });
        this.worker.removeEventListener('message', handler);
      }
    };
    this.worker.addEventListener('message', handler);
    this.worker.postMessage({ type: 'LAYOUT_START', token, nodes, edges } satisfies LayoutWorkerRequest);
  }
}

import type { GraphNode, GraphEdge } from '../../state.js';

export interface LayoutResult { nodes: GraphNode[]; edges: GraphEdge[]; iterations: number; durationMs: number }

export class LayoutWorkerAdapter {
  private worker: Worker;
  private currentToken = 0;
  constructor(workerPath: string) { this.worker = new Worker(workerPath, { type: 'module' }); }
  run(nodes: GraphNode[], edges: GraphEdge[], onDone: (r: LayoutResult) => void): void {
    this.currentToken += 1;
    const token = this.currentToken;
    const handler = (ev: MessageEvent): void => {
      if (ev.data?.token !== token || ev.data?.type !== 'layoutDone') return;
      this.worker.removeEventListener('message', handler);
      onDone({ nodes: ev.data.nodes, edges: ev.data.edges, iterations: ev.data.iterations, durationMs: ev.data.durationMs });
    };
    this.worker.addEventListener('message', handler);
    this.worker.postMessage({ token, nodes, edges });
  }
}

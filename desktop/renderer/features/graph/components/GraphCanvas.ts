import { GraphWebGLRenderer } from '../engine/renderer/webgl.js';
import { pickNode } from '../engine/picking/pick.js';
import type { GraphState } from '../state.js';

export function mountGraphCanvas(stateRef: () => GraphState, onSelect: (id: string) => void, onFocus: (id: string) => void): { render: () => void; setViewport: (zoom: number, panX: number, panY: number) => void } {
  const canvas = document.getElementById('graph-canvas') as HTMLCanvasElement;
  canvas.width = canvas.clientWidth || 640;
  canvas.height = canvas.clientHeight || 420;
  const renderer = new GraphWebGLRenderer(canvas);
  let zoom = 1; let panX = 0; let panY = 0;
  let dragging = false; let sx = 0; let sy = 0;
  canvas.onmousedown = (e) => { dragging = true; sx = e.clientX; sy = e.clientY; };
  canvas.onmouseup = (e) => {
    dragging = false;
    const s = stateRef();
    const id = pickNode(s.dataset.nodes, e.offsetX, e.offsetY, zoom, panX, panY, canvas.width, canvas.height);
    if (id) onSelect(id);
  };
  canvas.ondblclick = (e) => {
    const s = stateRef();
    const id = pickNode(s.dataset.nodes, e.offsetX, e.offsetY, zoom, panX, panY, canvas.width, canvas.height);
    if (id) onFocus(id);
  };
  canvas.onmousemove = (e) => { if (!dragging) return; panX += e.clientX - sx; panY += e.clientY - sy; sx = e.clientX; sy = e.clientY; };
  canvas.onwheel = (e) => { e.preventDefault(); zoom = Math.max(0.2, Math.min(4, zoom + (e.deltaY < 0 ? 0.1 : -0.1))); };
  return {
    render: () => { const s = stateRef(); renderer.render(s.dataset.nodes, s.dataset.edges, zoom, panX, panY); },
    setViewport: (z, x, y) => { zoom = z; panX = x; panY = y; }
  };
}

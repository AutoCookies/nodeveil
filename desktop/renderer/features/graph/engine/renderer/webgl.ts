import type { GraphEdge, GraphNode } from '../../state.js';

export class GraphWebGLRenderer {
  private gl: WebGLRenderingContext;
  constructor(private canvas: HTMLCanvasElement) {
    const gl = canvas.getContext('webgl');
    if (!gl) throw new Error('WebGL unavailable');
    this.gl = gl;
  }
  render(nodes: GraphNode[], edges: GraphEdge[], zoom: number, panX: number, panY: number): void {
    const gl = this.gl;
    gl.viewport(0, 0, this.canvas.width, this.canvas.height);
    gl.clearColor(0.98, 0.98, 0.99, 1);
    gl.clear(gl.COLOR_BUFFER_BIT);
    // fallback immediate mode emulation: use 2d overlay for now via webgl clear + cheap draw markers with scissor
    // keep CPU small and render loop explicit
    const ctx2d = this.canvas.getContext('2d');
    if (!ctx2d) return;
    ctx2d.clearRect(0, 0, this.canvas.width, this.canvas.height);
    ctx2d.strokeStyle = '#b5b5c7';
    edges.forEach((e) => {
      const a = nodes.find((n) => n.id === e.from_id); const b = nodes.find((n) => n.id === e.to_id);
      if (!a || !b) return;
      ctx2d.beginPath();
      ctx2d.moveTo(a.x * zoom + panX + this.canvas.width / 2, a.y * zoom + panY + this.canvas.height / 2);
      ctx2d.lineTo(b.x * zoom + panX + this.canvas.width / 2, b.y * zoom + panY + this.canvas.height / 2);
      ctx2d.stroke();
    });
    nodes.forEach((n) => {
      ctx2d.fillStyle = '#2f6fed';
      ctx2d.beginPath();
      ctx2d.arc(n.x * zoom + panX + this.canvas.width / 2, n.y * zoom + panY + this.canvas.height / 2, 3, 0, Math.PI * 2);
      ctx2d.fill();
    });
  }
}

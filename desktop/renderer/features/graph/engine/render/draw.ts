import type { GraphBuffers } from './buffers.js';
export function drawGraph(gl: WebGLRenderingContext, buffers: GraphBuffers): void {
  gl.clearColor(0.97, 0.97, 0.98, 1);
  gl.clear(gl.COLOR_BUFFER_BIT);
  void buffers;
}

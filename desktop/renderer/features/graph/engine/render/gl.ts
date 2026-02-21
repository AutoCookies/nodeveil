export function createGL(canvas: HTMLCanvasElement): WebGLRenderingContext {
  const gl = canvas.getContext('webgl');
  if (!gl) throw new Error('WebGL unavailable');
  return gl;
}

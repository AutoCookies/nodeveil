import { Edge, IPCError, IndexStatus, SearchResult, VersionInfo } from '../shared/contracts.js';
const ENGINE_URL = process.env.NODEVEIL_ENGINE_URL ?? 'http://127.0.0.1:8080';
async function req<T>(path: string, init?: RequestInit): Promise<T> {
  try { const r = await fetch(`${ENGINE_URL}${path}`, init); if (!r.ok) throw new IPCError(`status ${r.status}`); return (await r.json()) as T; }
  catch (e) { throw new IPCError(`request failed: ${path}`, e); }
}
export const engineClient = {
  getVersion: () => req<VersionInfo>('/version'),
  getStatus: () => req<IndexStatus>('/index/status'),
  listRoots: async () => (await req<{roots:string[]}>('/roots')).roots,
  addRoot: async (path: string) => { await req('/roots', { method:'POST', headers:{'content-type':'application/json'}, body: JSON.stringify({path}) }); },
  removeRoot: async (path: string) => { await req(`/roots?path=${encodeURIComponent(path)}`, { method: 'DELETE' }); },
  search: async (q: string) => (await req<{results:SearchResult[]}>(`/search?q=${encodeURIComponent(q)}`)).results,
  createLink: async (fromId: string, toId: string, relationType = 'related') => { await req('/graph/links', { method:'POST', headers:{'content-type':'application/json'}, body: JSON.stringify({fromId,toId,relationType}) }); },
  listLinks: async (nodeId: string, direction = 'BOTH') => (await req<{links:Edge[]}>(`/graph/links?nodeId=${encodeURIComponent(nodeId)}&direction=${direction}`)).links
};

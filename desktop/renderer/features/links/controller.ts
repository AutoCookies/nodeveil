import { requestJSON } from '../../shared/ipc/client.js';
import type { LinkItem } from './state.js';
export async function loadOutgoing(nodeId: string): Promise<LinkItem[]> {
  return (await requestJSON<{ links: LinkItem[] }>(`/graph/links?nodeId=${encodeURIComponent(nodeId)}&direction=OUT`)).links;
}
export async function loadIncoming(nodeId: string): Promise<LinkItem[]> {
  return (await requestJSON<{ links: LinkItem[] }>(`/graph/backlinks?nodeId=${encodeURIComponent(nodeId)}`)).links;
}
export async function createLink(fromId: string, toId: string, relationType: string, note = ''): Promise<void> {
  await requestJSON('/graph/links', { method: 'POST', body: { fromId, toId, relationType, note } });
}
export async function removeLink(id: string): Promise<void> {
  await requestJSON(`/graph/links?id=${encodeURIComponent(id)}`, { method: 'DELETE' });
}

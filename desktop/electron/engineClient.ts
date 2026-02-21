import { EngineClient, HealthStatus, IPCError, VersionInfo } from '../shared/contracts.js';

const ENGINE_URL = process.env.NODEVEIL_ENGINE_URL ?? 'http://127.0.0.1:8080';

async function getJSON<T>(path: string): Promise<T> {
  try {
    const response = await fetch(`${ENGINE_URL}${path}`);
    if (!response.ok) {
      throw new IPCError(`Engine request failed with status ${response.status}`);
    }
    return (await response.json()) as T;
  } catch (error) {
    throw new IPCError(`Failed to call engine path ${path}`, error);
  }
}

export const engineClient: EngineClient = {
  getVersion() {
    return getJSON<VersionInfo>('/version');
  },
  healthCheck() {
    return getJSON<HealthStatus>('/health');
  },
  async listRoots() {
    const response = await getJSON<{ roots: string[] }>('/roots');
    return response.roots;
  }
};

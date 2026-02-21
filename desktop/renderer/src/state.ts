export interface EngineState {
  statusText: string;
}

export function createInitialState(): EngineState {
  return { statusText: 'Engine status: checking...' };
}

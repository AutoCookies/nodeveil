export interface RetryPolicy { retries: number; backoffMs: number }
export interface RequestOptions {
  timeoutMs?: number;
  retry?: RetryPolicy;
  signal?: AbortSignal;
  method?: string;
  body?: unknown;
}

const BASE_URL = 'http://127.0.0.1:8080';

function sleep(ms: number): Promise<void> { return new Promise((r) => setTimeout(r, ms)); }

export async function requestJSON<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { timeoutMs = 4000, retry = { retries: 0, backoffMs: 100 }, signal, method = 'GET', body } = opts;
  let attempt = 0;
  while (true) {
    const ctl = new AbortController();
    const timer = setTimeout(() => ctl.abort(), timeoutMs);
    const onAbort = (): void => ctl.abort();
    signal?.addEventListener('abort', onAbort, { once: true });
    try {
      const init: RequestInit = { method, signal: ctl.signal };
      if (body) {
        init.headers = { 'content-type': 'application/json' };
        init.body = JSON.stringify(body);
      }
      const response = await fetch(`${BASE_URL}${path}`, init);
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      return (await response.json()) as T;
    } catch (error) {
      if (attempt >= retry.retries || method !== 'GET') throw error;
      attempt += 1;
      await sleep(retry.backoffMs * attempt);
    } finally {
      clearTimeout(timer);
      signal?.removeEventListener('abort', onAbort);
    }
  }
}

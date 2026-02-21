import test from 'node:test';
import assert from 'node:assert/strict';

async function requestJSON(path, opts = {}) {
  const { timeoutMs = 100, retry = { retries: 0, backoffMs: 1 }, method = 'GET' } = opts;
  let attempt = 0;
  while (true) {
    const ctl = new AbortController();
    const timer = setTimeout(() => ctl.abort(), timeoutMs);
    try {
      const r = await fetch(path, { signal: ctl.signal, method });
      if (!r.ok) throw new Error(`HTTP ${r.status}`);
      return await r.json();
    } catch (e) {
      if (attempt >= retry.retries || method !== 'GET') throw e;
      attempt += 1;
    } finally { clearTimeout(timer); }
  }
}

test('ipc client retries safe reads', async () => {
  let calls = 0;
  global.fetch = async () => { calls += 1; if (calls < 3) throw new Error('net'); return { ok: true, json: async () => ({ ok: true }) }; };
  const out = await requestJSON('http://x', { retry: { retries: 2, backoffMs: 1 } });
  assert.equal(out.ok, true);
  assert.equal(calls, 3);
});

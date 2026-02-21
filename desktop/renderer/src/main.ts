const API = 'http://127.0.0.1:8080';
async function j<T>(p: string, init?: RequestInit): Promise<T> { const r = await fetch(API+p, init); return (await r.json()) as T; }
async function refresh(): Promise<void> {
  const s = await j<any>('/index/status');
  (document.getElementById('status') as HTMLElement).textContent = `files=${s.filesIndexed} queue=${s.queueDepth} progress=${s.progressPercent}%`;
  (document.getElementById('errors') as HTMLElement).textContent = `errors=${s.errors} lastScan=${s.lastScanUnix}`;
  const roots = (await j<any>('/roots')).roots as string[];
  const ul = document.getElementById('roots') as HTMLUListElement; ul.innerHTML='';
  roots.forEach((r) => { const li = document.createElement('li'); li.textContent=r; ul.appendChild(li); });
}
(document.getElementById('add-root') as HTMLButtonElement).onclick = async () => {
  const path = (document.getElementById('root-input') as HTMLInputElement).value;
  await j('/roots', { method:'POST', headers:{'content-type':'application/json'}, body: JSON.stringify({path}) }); await refresh();
};
(document.getElementById('search-btn') as HTMLButtonElement).onclick = async () => {
  const q = (document.getElementById('search-input') as HTMLInputElement).value;
  const results = (await j<any>(`/search?q=${encodeURIComponent(q)}`)).results as Array<{absPath:string}>;
  const ul = document.getElementById('results') as HTMLUListElement; ul.innerHTML='';
  results.forEach((r) => { const li = document.createElement('li'); li.textContent=r.absPath; ul.appendChild(li); });
};
setInterval(() => { void refresh(); }, 1500);
void refresh();

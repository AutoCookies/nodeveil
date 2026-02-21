const API = 'http://127.0.0.1:8080';
async function j<T>(p: string, init?: RequestInit): Promise<T> { const r = await fetch(API+p, init); return (await r.json()) as T; }
async function refresh(): Promise<void> {
  const s = await j<any>('/index/status');
  (document.getElementById('status') as HTMLElement).textContent = `files=${s.filesIndexed} queue=${s.queueDepth} progress=${s.progressPercent}%`;
  (document.getElementById('errors') as HTMLElement).textContent = `errors=${s.errors} lastScan=${s.lastScanUnix}`;
}
(document.getElementById('add-root') as HTMLButtonElement).onclick = async () => {
  const path = (document.getElementById('root-input') as HTMLInputElement).value;
  await j('/roots', { method:'POST', headers:{'content-type':'application/json'}, body: JSON.stringify({path}) });
};
(document.getElementById('search-btn') as HTMLButtonElement).onclick = async () => {
  const q = (document.getElementById('search-input') as HTMLInputElement).value;
  const results = (await j<any>(`/search?q=${encodeURIComponent(q)}`)).results as Array<{id:string;absPath:string}>;
  const ul = document.getElementById('results') as HTMLUListElement; ul.innerHTML='';
  results.forEach((r) => { const li = document.createElement('li'); li.textContent=`${r.absPath} (${r.id})`; li.onclick=()=>{(document.getElementById('from-id') as HTMLInputElement).value=r.id;}; ul.appendChild(li); });
};
async function loadLinks(nodeId: string): Promise<void> {
  const out = (await j<any>(`/graph/links?nodeId=${encodeURIComponent(nodeId)}&direction=OUT`)).links as Array<any>;
  const inc = (await j<any>(`/graph/backlinks?nodeId=${encodeURIComponent(nodeId)}`)).links as Array<any>;
  const oul = document.getElementById('out-links') as HTMLUListElement; oul.innerHTML='';
  out.forEach((e) => { const li = document.createElement('li'); li.textContent = `${e.to_id} [${e.relation_type}]`; const b=document.createElement('button'); b.textContent='unlink'; b.onclick=async()=>{await j(`/graph/links?id=${encodeURIComponent(e.id)}`,{method:'DELETE'}); await loadLinks(nodeId);}; li.appendChild(b); oul.appendChild(li); });
  const iul = document.getElementById('in-links') as HTMLUListElement; iul.innerHTML='';
  inc.forEach((e) => { const li = document.createElement('li'); li.textContent = `${e.from_id} [${e.relation_type}]`; iul.appendChild(li); });
}
(document.getElementById('link-btn') as HTMLButtonElement).onclick = async () => {
  const fromId = (document.getElementById('from-id') as HTMLInputElement).value;
  const toId = (document.getElementById('to-id') as HTMLInputElement).value;
  const relationType = (document.getElementById('relation') as HTMLSelectElement).value;
  const error = document.getElementById('link-error') as HTMLElement;
  error.textContent = '';
  try { await j('/graph/links', { method:'POST', headers:{'content-type':'application/json'}, body: JSON.stringify({fromId,toId,relationType}) }); await loadLinks(fromId); }
  catch { error.textContent = 'Link failed (node not indexed, duplicate, or engine unavailable).'; }
};
setInterval(() => { void refresh(); }, 1500);
void refresh();

export function renderLegend(): void {
  const el = document.getElementById('graph-legend');
  if (!el) return;
  el.innerHTML = 'Relations: <span>related</span> <span>ref</span> <span>evidence</span> <span>depends_on</span>';
}

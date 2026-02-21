import type { GraphFilters } from '../state.js';
export function bindGraphToolbar(onApply: (f: GraphFilters) => void, onReset: () => void): void {
  (document.getElementById('graph-apply') as HTMLButtonElement).onclick = () => {
    const depth = Number((document.getElementById('graph-depth') as HTMLSelectElement).value) as 1 | 2 | 3;
    const direction = (document.getElementById('graph-direction') as HTMLSelectElement).value as GraphFilters['direction'];
    const rel = (document.getElementById('graph-rel') as HTMLInputElement).value.split(',').map((v) => v.trim()).filter(Boolean);
    const ext = (document.getElementById('graph-ext') as HTMLInputElement).value.split(',').map((v) => v.trim()).filter(Boolean);
    onApply({ depth, direction, relationTypes: rel, extensions: ext });
  };
  (document.getElementById('graph-reset') as HTMLButtonElement).onclick = onReset;
}

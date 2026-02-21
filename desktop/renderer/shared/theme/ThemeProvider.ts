import { loadTheme, saveTheme, themes, type ThemeName } from './useTheme.js';

export function applyTheme(name: ThemeName): void {
  const t = themes[name];
  const root = document.documentElement;
  Object.entries(t).forEach(([k, v]) => root.style.setProperty(`--nv-${k}`, v));
  saveTheme(name);
}
export function initTheme(): ThemeName {
  const n = loadTheme();
  applyTheme(n);
  return n;
}

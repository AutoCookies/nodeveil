import type { ThemeTokens } from './tokens.js';
import { lightTheme } from './themes/light.js';
import { darkTheme } from './themes/dark.js';
import { neonTheme } from './themes/neon.js';
import { sakuraTheme } from './themes/sakura.js';

export type ThemeName = 'light'|'dark'|'neon'|'sakura';
export const themes: Record<ThemeName, ThemeTokens> = { light: lightTheme, dark: darkTheme, neon: neonTheme, sakura: sakuraTheme };

export function loadTheme(): ThemeName {
  const v = localStorage.getItem('nodeveil.theme') as ThemeName | null;
  return v && themes[v] ? v : 'light';
}
export function saveTheme(name: ThemeName): void { localStorage.setItem('nodeveil.theme', name); }

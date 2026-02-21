declare module 'electron' {
  export const app: {
    whenReady(): Promise<void>;
  };
  export class BrowserWindow {
    constructor(opts: unknown);
    loadFile(path: string): Promise<void>;
  }
  export const Menu: {
    buildFromTemplate(template: unknown[]): unknown;
    setApplicationMenu(menu: unknown): void;
  };
}

declare const process: {
  env: Record<string, string | undefined>;
};

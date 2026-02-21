import { app, BrowserWindow, Menu } from 'electron';

function createWindow(): void {
  const win = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      contextIsolation: true
    }
  });

  const menu = Menu.buildFromTemplate([
    {
      label: 'File',
      submenu: [{ label: 'Open Root', accelerator: 'CmdOrCtrl+O' }, { role: 'quit' }]
    },
    {
      label: 'View',
      submenu: [{ role: 'reload' }, { role: 'toggledevtools' }]
    }
  ]);

  Menu.setApplicationMenu(menu);
  win.loadFile('renderer/index.html').catch(console.error);
}

app.whenReady().then(createWindow).catch(console.error);

$ErrorActionPreference = 'Stop'
Set-Location (Join-Path $PSScriptRoot '..')
pnpm -C desktop install --frozen-lockfile
pnpm -C desktop build
Write-Output "Run: pnpm -C desktop exec electron-builder --linux AppImage --win nsis --mac dmg"

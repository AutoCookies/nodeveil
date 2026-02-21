$ErrorActionPreference = 'Stop'
Set-Location (Join-Path $PSScriptRoot '..')
New-Item -ItemType Directory -Force -Path build/engine | Out-Null
Set-Location engine
go build -trimpath -ldflags "-s -w" -o ../build/engine/nodeveil-engine.exe ./cmd/nodeveil-engine

$ErrorActionPreference = "Stop"

Set-Location ./frontend

npm ci
npm run build

Set-Location ../
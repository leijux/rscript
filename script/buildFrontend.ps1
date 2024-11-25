$ErrorActionPreference = "Stop"

Set-Location ./frontend

npm ci
npm build

Set-Location ../
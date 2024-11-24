$ErrorActionPreference = "Stop"
$configPath = "./wails.json"
$config = Get-Content $configPath | ConvertFrom-Json

$config.info.productVersion = "$env:VERSION"
$config.outputfilename = "$($config.name)_gui_$env:VERSION"

$config | ConvertTo-Json | Set-Content $configPath

wails build -platform "windows/amd64" -tags "gui" -upx -ldflags "-s -w -X github.com/leijux/rscript/internal/pkg/version.Version=$env:VERSION"
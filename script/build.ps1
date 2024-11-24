$ErrorActionPreference = "Stop"
# $DebugPreference = "Continue"
Write-Debug "VERSION=$env:VERSION"

#git tag -a $env:VERSION -m $env:VERSION
$build_version_path = "./build/$env:VERSION"

if (-not (Test-Path -Path $build_version_path)) {
    New-Item -Path $build_version_path -ItemType Directory | Out-Null
}

gox -osarch "windows/amd64 linux/amd64 darwin/amd64" `
    -ldflags "-s -w -X github.com/leijux/rscript/internal/pkg/version.Version=$env:VERSION" `
    -output "./build/$env:VERSION/{{.Dir}}_tui_{{.OS}}_{{.Arch}}_$env:VERSION"

if (-not $?) {
    Write-Warning "gox build fail"
}

Get-ChildItem "./build/$env:VERSION/*" | ForEach-Object {
    upx -9 $_.FullName
}
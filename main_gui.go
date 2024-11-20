//go:build gui

package main

import (
	"embed"

	"github.com/leijux/rscript/internal/app/gui"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	err := gui.Main(assets)
	if err != nil {
		println("Error:", err.Error())
	}
}

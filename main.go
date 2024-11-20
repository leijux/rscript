//go:build !gui

package main

import (
	"github.com/leijux/rscript/internal/app/tui"
)

func main() {
	tui.Main()
}

package main

//go:generate go run ./version_gen.go opentracer

import (
	_ "embed"
	"github.com/davidalpert/opentracer/internal/cmd"
)

func main() {
	cmd.Execute()
}

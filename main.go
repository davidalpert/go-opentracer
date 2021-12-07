package main

//go:generate go run ./version_gen.go gopentracer

import (
	_ "embed"
	"github.com/davidalpert/gopentracer/internal/cmd"
)

func main() {
	cmd.Execute()
}

package main

import (
	_ "embed"
	"github.com/davidalpert/opentracer/internal/cmd"
)

func main() {
	cmd.Execute()
}

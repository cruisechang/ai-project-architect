package main

import (
	"fmt"
	"os"

	"project-generator/apa"
)

func main() {
	if err := apa.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

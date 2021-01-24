package main

import (
	"fmt"
	"os"

	"github.com/byxorna/regtest/pkg/cli"
)

func main() {
	app := cli.New(os.Args[1:])
	if err := app.Run(); err != nil {
		fmt.Printf("fuck: %v\n", err)
		os.Exit(1)
	}
}

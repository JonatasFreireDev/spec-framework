package main

import (
	"os"

	"github.com/JonatasFreireDev/spec-framework/internal/cli"
)

var version = "dev"

func main() {
	os.Exit(cli.New(version).Run(os.Args[1:], os.Stdout, os.Stderr))
}

package main

import (
	"fmt"
	"os"

	"github.com/joeldotdias/gat/internal/ops"
)

const HELP_TEXT = `usage: gat <command> [<args>...]
`

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, HELP_TEXT) //"usage: gat <command> [<args>...]\n")
		os.Exit(1)
	}

	currDir, _ := os.Getwd()
	repo := ops.Repo(currDir)

	switch command := args[1]; command {
	case "init":
		repo.Init()
	}
}

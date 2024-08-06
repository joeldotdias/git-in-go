package main

import (
	"fmt"
	"os"

	"github.com/joeldotdias/gat/internal/ops"
)

// TODO: doc all the commands
const HELP_TEXT = `usage: gat <command> [<args>...]
`

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, HELP_TEXT)
		os.Exit(1)
	}

	currDir, _ := os.Getwd()
	repo := ops.Repo(currDir)

	switch command := args[1]; command {
	case "init":
		repo.Init()
	case "cat-file":
		repo.CatFile(args[2], args[3])
	case "hash-object":
		// add some parsing logic here
		objType := "blob"
		write := false
		repo.HashObject(write, objType, args[2])
	case "ls-tree":
		recursive := false
		if len(args) > 3 && args[3] == "-r" {
			recursive = true
		}
		err := repo.TopLsTree(args[2], recursive)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}
}

package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/vderyagin/dfm/dotfile"
	"github.com/vderyagin/dfm/repo"
)

// Repo returns a repo.Repo object based on command line arguments
func Repo(c *cli.Context) *repo.Repo {
	return repo.New(
		c.GlobalString("store"),
		c.GlobalString("home"),
	)
}

// EnsureArgsPresent fails if no command line arguments provided
func EnsureArgsPresent(c *cli.Context) {
	if !c.Args().Present() {
		fmt.Fprintln(os.Stderr, "No arguments provided")
		os.Exit(1)
	}
}

// ArgDotFiles returns a collection of DotFile objects made from command line
// arguments
func ArgDotFiles(c *cli.Context) []*dotfile.DotFile {
	EnsureArgsPresent(c)

	repo := Repo(c)
	dotfiles := make([]*dotfile.DotFile, len(c.Args()))

	for idx, arg := range c.Args() {
		orig, err := filepath.Abs(arg)

		if err != nil {
			log.Fatal(err)
		}

		dotfiles[idx] = dotfile.New(repo.StoredFilePath(orig), orig)
	}

	return dotfiles
}

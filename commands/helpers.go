package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/vderyagin/dfm/dotfile"
	"github.com/vderyagin/dfm/logger"
	"github.com/vderyagin/dfm/repo"
)

// Repo returns a repo.Repo object based on command line arguments.
func Repo(c *cli.Context) *repo.Repo {
	return repo.New(
		c.GlobalString("store"),
		c.GlobalString("home"),
	)
}

// EnsureArgsPresent fails loudly if no command line arguments provided.
func EnsureArgsPresent(c *cli.Context) {
	if !c.Args().Present() {
		fmt.Fprintln(os.Stderr, "No arguments provided")
		os.Exit(1)
	}
}

// ArgDotFiles returns a collection of DotFile objects constructed according
// to provided command line arguments.
func ArgDotFiles(c *cli.Context) []*dotfile.DotFile {
	EnsureArgsPresent(c)

	repo := Repo(c)
	dotfiles := make([]*dotfile.DotFile, len(c.Args()))

	for idx, arg := range c.Args() {
		orig, err := filepath.Abs(arg)

		if err != nil {
			log.Fatal(err)
		}

		if stored, err := repo.StoredFilePath(orig, c.Bool("host-specific"), c.Bool("copy")); err != nil {
			log.Fatal(err)
		} else {
			dotfiles[idx] = dotfile.New(stored, orig)
		}
	}

	return dotfiles
}

// Logger returns a Logger object for given dotfile.
func Logger(c *cli.Context, df *dotfile.DotFile) *logger.Logger {
	repo := Repo(c)
	id, _ := filepath.Rel(repo.Store, df.StoredLocation)
	return logger.New(id)
}

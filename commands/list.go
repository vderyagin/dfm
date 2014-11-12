package commands

import (
	"fmt"
	"path/filepath"

	"github.com/codegangsta/cli"
)

// List displays a list of stored dotfiles and their states.
func List(c *cli.Context) {
	for _, df := range Repo(c).StoredDotFiles() {
		id, _ := filepath.Rel(Repo(c).Store, df.StoredLocation)
		fmt.Printf("%23s %s\n", df.CurrentState().ColorString(), id)
	}
}

package commands

import (
	"fmt"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

// List displays a list of stored dotfiles and their states.
func List(c *cli.Context) error {
	for df := range Repo(c).StoredDotFiles() {
		id, _ := filepath.Rel(Repo(c).Store, df.StoredLocation)
		fmt.Printf("%23s %s\n", df.CurrentState().ColorString(), id)
	}

	return nil
}

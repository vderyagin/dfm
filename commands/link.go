package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Link links all stored dotfiles to their respective places in home
// directory.
func Link(c *cli.Context) {
	for _, df := range Repo(c).StoredDotFiles() {
		if err := df.Link(); err != nil {
			fmt.Println(err)
		}
	}
}

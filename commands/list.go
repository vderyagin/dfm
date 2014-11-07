package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// List displays a list of stored dotfiles.
func List(c *cli.Context) {
	for _, df := range Repo(c).StoredDotFiles() {
		fmt.Println(df)
	}
}

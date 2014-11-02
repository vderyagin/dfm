package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/vderyagin/dfm/repo"
)

// List displays a list of stored dotfiles.
func List(c *cli.Context) {
	store := c.GlobalString("store")
	home := c.GlobalString("home")

	for _, df := range repo.New(store, home).StoredDotFiles() {
		fmt.Println(df)
	}
}

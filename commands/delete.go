package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Delete removes both stored file and it's symlink.
func Delete(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		if err := df.Delete(); err != nil {
			fmt.Println(err)
		}
	}
}

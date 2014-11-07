package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Store stores and links back given files.
func Store(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		if err := df.Store(); err != nil {
			fmt.Println(err)
		}
	}
}

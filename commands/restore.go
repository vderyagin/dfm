package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// Restore moves file from store back to its original location.
func Restore(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		if err := df.Restore(); err != nil {
			fmt.Println(err)
		}
	}
}

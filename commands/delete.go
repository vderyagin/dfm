package commands

import "github.com/codegangsta/cli"

// Delete removes both stored file and it's symlink.
func Delete(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)

		if err := df.Delete(); err == nil {
			logger.Success("deleted")
		} else {
			logger.Fail("failed to delete", err.Error())
		}
	}
}

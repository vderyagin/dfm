package commands

import "github.com/codegangsta/cli"

// Restore moves file from store back to its original location.
func Restore(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)
		if err := df.Restore(); err == nil {
			logger.Success("restored")
		} else {
			logger.Fail("failed to restore", err.Error())
		}
	}
}

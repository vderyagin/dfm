package commands

import "github.com/codegangsta/cli"

// Store stores and links back given files.
func Store(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)
		if err := df.Store(); err == nil {
			logger.Success("stored")
		} else {
			logger.Fail("failed to store", err.Error())
		}
	}
}

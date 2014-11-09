package commands

import "github.com/codegangsta/cli"

// Link links all stored dotfiles to their respective places in home
// directory.
func Link(c *cli.Context) {
	for _, df := range Repo(c).StoredDotFiles() {
		if df.IsLinked() {
			continue
		}

		logger := Logger(c, df)

		if err := df.Link(); err == nil {
			logger.Success("linked")
		} else {
			logger.Fail("failed to link", err.Error())
		}
	}
}

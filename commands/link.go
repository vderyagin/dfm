package commands

import "github.com/codegangsta/cli"

// Link links all stored dotfiles to their respective places in home
// directory.
func Link(c *cli.Context) {
	for _, df := range Repo(c).StoredDotFiles() {
		logger := Logger(c, df)

		if df.IsLinked() {
			logger.Skip("skipped linking", "is already linked")
		} else if err := df.Link(); err == nil {
			logger.Success("linked")
		} else {
			logger.Fail("failed to link", err.Error())
		}
	}
}

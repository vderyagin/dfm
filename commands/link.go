package commands

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

// Link links all stored dotfiles to their respective locations in home
// directory.
func Link(c *cli.Context) error {
	for df := range Repo(c).StoredDotFiles() {
		if df.IsLinked() {
			continue
		}

		logger := Logger(c, df)

		if c.Bool("force") && df.IsStored() {
			if err := os.RemoveAll(df.OriginalLocation); err != nil {
				logger.Fail("failed to remove file", err.Error())
			}
		}

		if err := df.Link(); err == nil {
			logger.Success("linked")
		} else {
			logger.Fail("failed to link", err.Error())
		}
	}

	return nil
}

package commands

import (
	"os"

	"github.com/urfave/cli"
)

// Link links all stored dotfiles to their respective locations in home
// directory.
func Link(c *cli.Context) error {
	var errs []error

	for df := range Repo(c).StoredDotFiles() {
		if df.IsLinked() {
			continue
		}

		logger := Logger(c, df)

		if c.Bool("force") && df.IsStored() {
			if err := os.RemoveAll(df.OriginalLocation); err != nil {
				logger.Fail("failed to remove file", err.Error())
				errs = append(errs, err)
			}
		}

		if err := df.Link(); err == nil {
			logger.Success("linked")
		} else {
			logger.Fail("failed to link", err.Error())
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return cli.NewMultiError(errs...)
}

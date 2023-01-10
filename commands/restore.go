package commands

import (
	"github.com/urfave/cli"
	"github.com/vderyagin/dfm/dotfile"
)

// Restore moves dotfiles from store back to its original location, makes
// sense only for linked files.
func Restore(c *cli.Context) error {
	var errs []error

	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)

		err := df.Restore()

		if err != nil {
			errs = append(errs, err)
		}

		switch err.(type) {
		case nil:
			logger.Success("restored")
		case dotfile.SkipError:
			logger.Skip("skipped restoring", err.Error())
		default:
			logger.Fail("failed to restore", err.Error())
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return cli.NewMultiError(errs...)
}

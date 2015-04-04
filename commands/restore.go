package commands

import (
	"github.com/codegangsta/cli"
	"gitlab.com/vderyagin/dfm/dotfile"
)

// Restore moves dotfiles from store back to its original location, makes
// sense only for linked files.
func Restore(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)

		switch err := df.Restore().(type) {
		case nil:
			logger.Success("restored")
		case dotfile.SkipError:
			logger.Skip("skipped restoring", err.Error())
		default:
			logger.Fail("failed to restore", err.Error())
		}
	}
}

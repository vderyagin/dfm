package commands

import (
	"gitlab.com/vderyagin/dfm/dotfile"
	"gopkg.in/urfave/cli.v1"
)

// Delete removes both stored file and it's symlink, works for properly linked
// files only.
func Delete(c *cli.Context) error {
	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)

		switch err := df.Delete().(type) {
		case nil:
			logger.Success("deleted")
		case dotfile.SkipError:
			logger.Skip("skipped deleting", err.Error())
		default:
			logger.Fail("failed to delete", err.Error())
		}
	}

	return nil
}

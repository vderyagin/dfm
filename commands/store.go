package commands

import (
	"os"

	"github.com/codegangsta/cli"
	"gitlab.com/vderyagin/dfm/dotfile"
	"gitlab.com/vderyagin/dfm/fsutil"
)

// Store stores and links back given files.
func Store(c *cli.Context) error {
	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)

		if c.Bool("force") && fsutil.IsRegularFile(df.OriginalLocation) {
			if err := os.RemoveAll(df.StoredLocation); err != nil {
				logger.Fail("failed to remove file", err.Error())
			}
		}

		switch err := df.Store().(type) {
		case nil:
			logger.Success("stored")
		case dotfile.SkipError:
			logger.Skip("skipped storing", err.Error())
		default:
			logger.Fail("failed to store", err.Error())
		}
	}

	return nil
}

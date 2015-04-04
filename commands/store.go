package commands

import (
	"os"

	"github.com/codegangsta/cli"
	"gitlab.com/vderyagin/dfm/dotfile"
	"gitlab.com/vderyagin/dfm/fsutil"
)

// Store stores and links back given files.
func Store(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		if c.Bool("force") && fsutil.IsRegularFile(df.OriginalLocation) {
			os.RemoveAll(df.StoredLocation)
		}

		logger := Logger(c, df)

		switch err := df.Store().(type) {
		case nil:
			logger.Success("stored")
		case dotfile.SkipError:
			logger.Skip("skipped storing", err.Error())
		default:
			logger.Fail("failed to store", err.Error())
		}
	}
}

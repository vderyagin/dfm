package commands

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/vderyagin/dfm/fsutil"
)

// Store stores and links back given files.
func Store(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		if c.Bool("force") && fsutil.IsRegularFile(df.OriginalLocation) {
			os.RemoveAll(df.StoredLocation)
		}

		logger := Logger(c, df)

		if df.IsLinked() {
			logger.Skip("skipped storing", "is already stored and linked")
		} else if err := df.Store(); err == nil {
			logger.Success("stored")
		} else {
			logger.Fail("failed to store", err.Error())
		}
	}
}

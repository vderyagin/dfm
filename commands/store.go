package commands

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/vderyagin/dfm/dotfile"
	"github.com/vderyagin/dfm/fsutil"
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
		case dotfile.FailError:
			logger.Fail("failed to store", err.Error())
		default:
			log.Fatalf("error of unknown type: %v", err)
		}
	}
}

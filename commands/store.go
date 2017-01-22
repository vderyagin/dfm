package commands

import (
	"os"

	"github.com/urfave/cli"

	"gitlab.com/vderyagin/dfm/dotfile"
	"gitlab.com/vderyagin/dfm/fsutil"
)

// Store stores and links back given files.
func Store(c *cli.Context) error {
	var errs []error

	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)

		if c.Bool("force") && fsutil.IsRegularFile(df.OriginalLocation) {
			if err := os.RemoveAll(df.StoredLocation); err != nil {
				logger.Fail("failed to remove file", err.Error())
				errs = append(errs, err)
			}
		}

		err := df.Store()

		if err != nil {
			errs = append(errs, err)
		}

		switch err.(type) {
		case nil:
			logger.Success("stored")
		case dotfile.SkipError:
			logger.Skip("skipped storing", err.Error())
		default:
			logger.Fail("failed to store", err.Error())
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return cli.NewMultiError(errs...)
}

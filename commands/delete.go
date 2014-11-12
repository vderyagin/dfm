package commands

import (
	"github.com/codegangsta/cli"
	"github.com/vderyagin/dfm/fsutil"
)

// Delete removes both stored file and it's symlink, works for properly linked
// files only.
func Delete(c *cli.Context) {
	for _, df := range ArgDotFiles(c) {
		logger := Logger(c, df)

		if !(fsutil.Exists(df.OriginalLocation) || fsutil.Exists(df.StoredLocation)) {
			logger.Skip("skipped deleting", "does not exist")
		} else if err := df.Delete(); err == nil {
			logger.Success("deleted")
		} else {
			logger.Fail("failed to delete", err.Error())
		}
	}
}

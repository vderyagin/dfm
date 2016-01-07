package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/codegangsta/cli"
	"gitlab.com/vderyagin/dfm/commands"
)

func homeDir() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir
}

var appFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "home",
		Value:  homeDir(),
		Usage:  "home directory",
		EnvVar: "DOTFILES_HOME_DIR",
	},
	cli.StringFlag{
		Name:   "store",
		Value:  filepath.Join(homeDir(), ".dotfiles"),
		Usage:  "directory files will be stored in",
		EnvVar: "DOTFILES_STORE_DIR",
	},
}

var appCommands = []cli.Command{
	{
		Name:      "list",
		ShortName: "l",
		Usage:     "List stored dotfiles",
		Action:    commands.List,
	},
	{
		Name:      "store",
		ShortName: "s",
		Usage:     "Put file into store",
		Action:    commands.Store,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "force",
				Usage: "overwrite conflicting files if necessary",
			},
			cli.BoolFlag{
				Name:  "host-specific",
				Usage: "store for this host only (hosts are distinguished by hostnames)",
			},
			cli.BoolFlag{
				Name:  "copy",
				Usage: "make sure this file always gets copied, not symlinked",
			},
		},
	},
	{
		Name:      "restore",
		ShortName: "r",
		Usage:     "Move file to its original location",
		Action:    commands.Restore,
	},
	{
		Name:   "link",
		Usage:  "Link all stored files to their original locations",
		Action: commands.Link,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "force",
				Usage: "overwrite conflicting files if necessary",
			},
		},
	},
	{
		Name:      "delete",
		ShortName: "d",
		Usage:     "Delete given files from home and store",
		Action:    commands.Delete,
	},
}

func main() {
	app := cli.NewApp()

	app.Author = "Victor Deryagin <vderyagin@gmail.com>"
	app.Name = "dfm"
	app.Usage = "Dotfile manager"
	app.Version = "0.2.0"
	app.Flags = appFlags
	app.Commands = appCommands

	app.Run(os.Args)
}

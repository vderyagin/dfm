package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/vderyagin/dfm/commands"
)

func homeDir() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir
}

func main() {
	app := cli.NewApp()

	app.Author = "Victor Deryagin <vderyagin@gmail.com>"
	app.Name = "dfm"
	app.Usage = "Dotfile manager"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "home",
			Value:  homeDir(),
			Usage:  "Home directory",
			EnvVar: "DOTFILES_HOME_DIR",
		},
		cli.StringFlag{
			Name:   "store",
			Value:  filepath.Join(homeDir(), "dotfiles-test"),
			Usage:  "directory files will be stored in",
			EnvVar: "DOTFILES_STORE_DIR",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "List stored dotfiles",
			Action:    commands.List,
		},
	}

	app.Run(os.Args)
}

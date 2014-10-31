package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/vderyagin/dfm/commands"
)

func main() {
	app := cli.NewApp()

	app.Author = "Victor Deryagin <vderyagin@gmail.com>"
	app.Name = "dfm"
	app.Usage = "Dotfile manager"
	app.Version = "0.0.1"

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

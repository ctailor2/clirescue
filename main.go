package main

import (
	"net/http"
	"bufio"
	"os"
	u "os/user"

	"github.com/ctailor2/clirescue/trackerapi"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "clirescue"
	app.Usage = "CLI tool to talk to the Pivotal Tracker's API"

	app.Commands = []cli.Command{
		{
			Name:  "me",
			Usage: "prints out Tracker's representation of your account",
			Action: func(c *cli.Context) {
				client := &http.Client{}
				usr, _ := u.Current()
				trackerapi.Me(os.Stdout, bufio.NewReader(os.Stdin), client, usr.HomeDir)
			},
		},
	}

	app.Run(os.Args)
}

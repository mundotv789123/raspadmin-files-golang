package cmd

import (
	ctx "context"
	"fmt"

	"github.com/urfave/cli/v3"
)

var webCommand = cli.Command{
	Name:   "web",
	Usage:  "Start the web server",
	Action: runWeb,
}

func runWeb(_ ctx.Context, cmd *cli.Command) error {
	fmt.Println("Starting web server...")
	return nil
}

package cmd

import (
	ctx "context"

	"github.com/gin-gonic/gin"
	raroute "github.com/mundotv789123/raspadmin/router"
	"github.com/urfave/cli/v3"
)

var webCommand = cli.Command{
	Name:   "web",
	Usage:  "Start the web server",
	Action: runWeb,
}

func runWeb(_ ctx.Context, cmd *cli.Command) error {
	router := gin.Default()
	router.GET("/api", raroute.Index)
	router.Run()
	return nil
}

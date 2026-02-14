package cmd

import (
	ctx "context"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/database"
	"github.com/mundotv789123/raspadmin/router"
	"github.com/urfave/cli/v3"
)

var webCommand = cli.Command{
	Name:   "web",
	Usage:  "Start the web server",
	Action: runWeb,
}

func runWeb(_ ctx.Context, cmd *cli.Command) error {
	_, err := database.OpenDbConnection()
	if err != nil {
		return err
	}

	r := gin.Default()

	webCtx := &router.WebContext{DB: database.DB}
	webCtx.Routers(r)

	return r.Run("0.0.0.0:8080")
}

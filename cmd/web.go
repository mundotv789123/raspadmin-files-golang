package cmd

import (
	ctx "context"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/database"
	icongenerator "github.com/mundotv789123/raspadmin/jobs/icon_generator"
	"github.com/mundotv789123/raspadmin/router"
	"github.com/urfave/cli/v3"
)

var webCommand = cli.Command{
	Name:   "web",
	Usage:  "Start the web server",
	Action: runWeb,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "gen-thumbnail",
			Usage: "Enable background service to generate file thumbnails",
		},
	},
}

func runWeb(_ ctx.Context, cmd *cli.Command) error {
	_, err := database.OpenDbConnection()
	if err != nil {
		return err
	}

	genThumb := cmd.Bool("gen-thumbnail")
	if genThumb {
		icongenerator.StartBackgroundService()
	}

	r := gin.Default()

	webCtx := &router.WebContext{DB: database.DB}
	webCtx.Routers(r)

	return r.Run("0.0.0.0:8080")
}


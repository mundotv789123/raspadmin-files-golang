package cmd

import (
	ctx "context"

	"github.com/gin-gonic/gin"
	"github.com/mundotv789123/raspadmin/internal/config"
	"github.com/mundotv789123/raspadmin/internal/database"
	"github.com/mundotv789123/raspadmin/internal/router"
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

	r.Use(corsMiddleware())

	apiRouter := r.Group("/api")
	apiRouter.GET("", router.Index)
	apiRouter.GET("files", router.Files)
	apiRouter.GET("files/open", router.OpenFile)

	r.Run()
	return nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isOriginAllowed := func(origin string, allowedOrigins []string) bool {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		}
		origin := c.Request.Header.Get("Origin")

		if isOriginAllowed(origin, config.AllowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

package router

import (
	"print-service/internal/config"
	"print-service/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, cfg *config.Config) {
	group := r.Group("") // jlaisse un grp vide pour le moment, pour les middleware etc, vu que j'apprends
	{
		group.POST("/print", func(c *gin.Context) {
			controller.HandlePrint(c, cfg)
		})
	}

}

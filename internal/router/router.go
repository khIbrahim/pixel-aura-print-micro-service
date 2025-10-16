package router

import (
	"log"
	"print-service/internal/api"
	"print-service/internal/config"
	"print-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, cfg *config.Config) {
	printService, err := service.NewPrintService(cfg)
	if err != nil {
		log.Fatalf("erreur lors de l'initialisation du service d'impression : %v\n", err)
	}

	validator := api.NewTicketValidator(*cfg)
	handler := api.NewPrintHandler(validator, printService)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/print", handler.HandlePrintRequest)
		//v1.GET("/jobs/:id", handler.HandleJobStatus)
	}
	
}

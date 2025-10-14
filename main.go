package main

import (
	"fmt"
	"log"
	"print-service/internal/config"
	"print-service/internal/logger"
	"print-service/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erreur chargement configuration : %v", err)
	}

	err = logger.Init(cfg)
	if err != nil {
		log.Fatalf("Erreur initialisation logger : %v", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	router.SetupRouter(r, cfg)
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	err = r.Run(port)
	if err != nil {
		logger.Log.Fatalf("erreur lors du démarrage du serveur: %v", err)
		return
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"print-service/internal/config"
	"print-service/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erreur chargement configuration : %v", err)
	}

	fmt.Printf("Configuration chargée : %+v\n", cfg)

	logFile, _ := os.Create("print-service.log")
	if logFile != nil {
		log.SetOutput(logFile)
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
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
		return
	}
}

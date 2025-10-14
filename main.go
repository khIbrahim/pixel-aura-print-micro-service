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
		log.Fatalf("erreur lors du chargement de la configuration : %v\n", err)
	}

	err = logger.Init(cfg)
	if err != nil {
		log.Fatalf("erreur lors de l'initialisation du logger : %v\n", err)
	}

	//ticketStorage, er := ticketfiles.NewTicketStorage(&cfg.Print)
	//if er != nil {
	//	log.Fatalf("erreur lors de l'initialisation du stockage des tickets\n")
	//}
	//
	//filePath, err := ticketStorage.CreateTempFile("-----------------------------------------\n         TICKET RESTAURANT\n-----------------------------------------\nÉtablissement : Le Bistrot du Coin\nAdresse        : 12 Rue de la République\n                 75001 Paris\nTél.           : 01 42 56 78 90\n-----------------------------------------\nDate : 14/10/2025      Heure : 12:47\nServeur : Julie\n-----------------------------------------\nQté   Désignation              Prix (€)\n-----------------------------------------\n1     Plat du jour             12.50\n1     Dessert maison            5.00\n1     Café                      2.00\n-----------------------------------------\n           Total HT :          16.67\n           TVA (10%) :          1.83\n-----------------------------------------\n           TOTAL TTC :        18.50 €\n-----------------------------------------\nMode de paiement : TICKET RESTAURANT\n-----------------------------------------\n      Merci et à bientôt !\n-----------------------------------------\n")
	//if err != nil {
	//	log.Fatalf("erreur lors de la création du fichier temporaire : %v\n", err)
	//}
	//
	//logger.Log.Infof("\nFichier temporaire créé avec succès : %s\n", filePath)

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

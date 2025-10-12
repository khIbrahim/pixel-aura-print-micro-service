package controller

import (
	"fmt"

	"print-service/internal/config"
	"print-service/internal/service"

	"github.com/gin-gonic/gin"
)

func HandlePrint(c *gin.Context, cfg *config.Config) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": "false",
			"message": "machi le bon format",
		})
		return
	}

	fmt.Printf("Requête reçue : %+v\n", req)

	content := req.Content + "\n"
	printService := service.NewPrintService(cfg)
	err = printService.PrintToFileAndSend(content)
	if err != nil {
		c.JSON(500, gin.H{
			"success": "false",
			"message": fmt.Sprintf("erreur impression : %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": "true",
		"message": "impression lancée",
	})
}

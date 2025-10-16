package api

import (
	"fmt"
	"print-service/internal/config"
	"print-service/internal/logger"
	"print-service/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type TicketValidator struct {
	config config.Config
}

func NewTicketValidator(cfg config.Config) *TicketValidator {
	return &TicketValidator{
		config: cfg,
	}
}

func (validator *TicketValidator) ValidateTicketRequest(c *gin.Context) (*models.TicketRequest, error) {
	var req models.TicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la requête: %w", err)
	}

	switch req.Type {
	case models.KitchenTicket, models.BarTicket, models.CashierTicket, models.CustomerTicket:
	//valide
	default:
		return nil, fmt.Errorf("le type %s n'est pas valide", req.Type)
	}

	if strings.TrimSpace(req.PrinterName) == "" {
		return nil, fmt.Errorf("le nom de l'imprimante ne peut pas être vide")
	}

	if req.Copies < 0 || req.Copies > 5 {
		return nil, fmt.Errorf("le nombre de copies doit être entre 1 et 5")
	}

	if req.Priority < models.LowPriority || req.Priority > models.HighPriority {
		return nil, fmt.Errorf("la priorité %d n'est pas valide", req.Priority)
	}

	logger.Log.Printf("Taille du contenu reçu : %d octets", len(req.Content))

	if len(req.Content) <= 0 || len(req.Content) > validator.config.Print.MaxContentSize {
		return nil, fmt.Errorf("la taille du contenu doit être entre 0 et %d octets", validator.config.Print.MaxContentSize)
	}

	logger.Log.WithFields(map[string]interface{}{
		"type":     req.Type,
		"printer":  req.PrinterName,
		"order_id": req.OrderID,
		"priority": req.Priority,
	}).Debug("Ticket Request Validé")

	return &req, nil
}

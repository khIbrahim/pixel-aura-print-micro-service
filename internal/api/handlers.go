package api

import (
	"net/http"
	"print-service/internal/logger"
	"print-service/internal/models"
	"print-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PrintHandler struct {
	validator *TicketValidator
	printer   *service.PrintService
}

func NewPrintHandler(validator *TicketValidator, printer *service.PrintService) *PrintHandler {
	return &PrintHandler{
		validator: validator,
		printer:   printer,
	}
}

func (handler *PrintHandler) HandlePrintRequest(c *gin.Context) {
	req, err := handler.validator.ValidateTicketRequest(c)
	if err != nil {
		logger.Log.WithError(err).Error("validation de la requête échouée")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Requête invalide",
			"error":   err.Error(),
		})
		return
	}

	job := models.NewTicketJob(req)
	logger.Log.WithField("job_id", job.ID).Info("nouveau job de ticket créé")

	go func() {
		if err = handler.printer.PrintTicketJob(job); err != nil {
			logger.Log.WithError(err).WithField("job_id", job.ID).Error("échec d'impression")
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"job_id":  job.ID,
		"message": "Job d'impression ajouté à la file d'attente",
	})
}

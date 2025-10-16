package service

import (
	"fmt"
	"os/exec"
	"print-service/internal/config"
	"print-service/internal/logger"
	"print-service/internal/models"
	"print-service/internal/ticketfiles"
	"sync"
)

type PrintService struct {
	cfg           *config.Config
	ticketStorage *ticketfiles.TicketStorage
	jobs          map[string]*models.TicketJob
	mu            sync.RWMutex
}

func NewPrintService(cfg *config.Config) (*PrintService, error) {
	ticketStorage, er := ticketfiles.NewTicketStorage(&cfg.Print)
	if er != nil {
		return nil, fmt.Errorf("erreur lors de l'initialisation du stockage des tickets\n")
	}

	ticketStorage.LunchCleanupTask()

	return &PrintService{
		cfg:           cfg,
		ticketStorage: ticketStorage,
		jobs:          make(map[string]*models.TicketJob),
	}, nil
}

func (service *PrintService) PrintTicketJob(job *models.TicketJob) error {
	service.mu.Lock()
	service.jobs[job.ID] = job
	service.mu.Unlock()

	job.UpdateStatus(models.PrintingStatus, nil)
	filePath, err := service.ticketStorage.CreateTempFile(job.Request.Content)
	if err != nil {
		job.UpdateStatus(models.FailedStatus, nil)
		return fmt.Errorf("erreur lors de la création du fichier temporaire : %v\n", err)
	}

	logger.Log.Infof("Fichier temporaire créé avec succès : %s\n", filePath)
	job.FilePath = filePath

	for i := 0; i < job.Copies; i++ {
		logger.Log.WithFields(map[string]interface{}{
			"job_id":  job.ID,
			"printer": job.PrinterName,
			"copy":    i,
			"total":   job.Copies,
		}).Debug("Impression d'une copie du ticket")

		cmd := fmt.Sprintf(service.cfg.Print.Command, filePath)
		psCmd := exec.Command("powershell", "-Command", cmd)

		if output, err := psCmd.CombinedOutput(); err != nil {
			logger.Log.WithError(err).WithField("output", string(output)).Error("Erreur lors de l'impression")
			job.UpdateStatus("failed", err)
			return fmt.Errorf("erreur d'impression: %w", err)
		}
	}

	if err := service.ticketStorage.ArchiveTicket(job); err != nil {
		logger.Log.WithError(err).Warn("erreur lors de l'archivage du ticket")
	}

	job.UpdateStatus(models.CompletedStatus, nil)

	logger.Log.WithFields(map[string]interface{}{
		"job_id":   job.ID,
		"printer":  job.PrinterName,
		"order_id": job.Request.OrderID,
		"type":     job.Request.Type,
		"copies":   job.Copies,
	}).Info("Job d'impression terminé avec succès")

	return nil
}

func (service *PrintService) GetJobStatus(jobID string) *models.TicketJob {
	service.mu.RLock()
	defer service.mu.RUnlock()

	return service.jobs[jobID]
}

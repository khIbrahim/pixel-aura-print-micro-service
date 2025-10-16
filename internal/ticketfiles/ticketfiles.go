package ticketfiles

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"print-service/internal/config"
	"print-service/internal/logger"
	"print-service/internal/models"
	"print-service/internal/tasks"
	"time"

	"github.com/google/uuid"
)

type TicketStorage struct {
	cfg         *config.PrintConfig
	tempDir     string
	archiveDir  string
	keepArchive bool
}

type TicketMetadata struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"order_id"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	PrintedAt   time.Time `json:"printed_at"`
	Content     string    `json:"content"`
	PrinterName string    `json:"printer_name"`
}

func CreateTicketMetadata(orderID, ticketType, content, printerName string, createdAt time.Time) TicketMetadata {
	return TicketMetadata{
		ID:          uuid.New().String(),
		OrderID:     orderID,
		Type:        ticketType,
		CreatedAt:   createdAt,
		Content:     content,
		PrinterName: printerName,
	}
}

func NewTicketStorage(config *config.PrintConfig) (*TicketStorage, error) {
	if err := os.MkdirAll(config.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("erreur lors de la création du dossier temporaire : %v\n", err)
	}

	archiveDir := config.ArchiveDir
	keepArchive := config.KeepArchive

	if config.KeepArchive {
		if err := os.MkdirAll(archiveDir, 0755); err != nil {
			return nil, fmt.Errorf("erreur lors de la création du dossier d'archive : %v\n", err)
		}
	}

	return &TicketStorage{
		config,
		config.TempDir,
		archiveDir,
		keepArchive,
	}, nil
}

func (storage *TicketStorage) CreateTempFile(content string) (string, error) {
	ticketID := uuid.New().String()
	fileName := "ticket_" + ticketID + ".txt"
	filePath := filepath.Join(storage.tempDir, fileName)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du fichier temporaire %s : %v\n", filePath, err)
	}

	logger.Log.WithField("path", filePath).Info("Fichier temporaire créé avec succès")

	return filePath, nil
}

func (storage *TicketStorage) ArchiveTicket(ticket *models.TicketJob) error {
	if !storage.keepArchive {
		return nil
	}

	now := time.Now()
	archivePath := filepath.Join(
		storage.archiveDir,
		fmt.Sprintf("%d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()),
	)

	if err := os.MkdirAll(archivePath, 0755); err != nil {
		return fmt.Errorf("erreur lors de la création du dossier d'archive %s : %v\n", archivePath, err)
	}

	fileName := fmt.Sprintf("order_%s_%s.json", ticket.Request.OrderID, ticket.Request.Type)
	filePath := filepath.Join(archivePath, fileName)

	data, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de la conversion de ticket en JSON : %v\n", err)
	}

	if err = os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("erreur lors de la création du fichier d'archive %s : %v\n", filePath, err)
	}

	logger.Log.WithField("path", filePath).Info("Ticket archivé avec succès")
	return nil
}

func (storage *TicketStorage) CleanupTempFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("erreur lors de la suppression du fichier temporaire %s : %v\n", filePath, err)
	}

	logger.Log.WithField("path", filePath).Info("Fichier temporaire supprimé avec succès")
	return nil
}

func (storage *TicketStorage) LunchCleanupTask() {
	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	cleanupInterval := time.Duration(storage.cfg.CleanupInterval)
	go tasks.CleanTempTickets(ctx, storage.cfg.TempDir, cleanupInterval*time.Minute)
}

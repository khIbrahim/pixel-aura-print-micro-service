package service

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"print-service/internal/config"
	"print-service/internal/logger"
	"print-service/internal/tasks"
	"print-service/internal/ticketfiles"
	"time"
)

type PrintService struct {
	cfg *config.Config
}

func NewPrintService(cfg *config.Config) *PrintService {
	// todo entourer le CleanTempTickets avec une sorte de try catch pour utiliser cancel et éviter que goroutine plante le service
	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	cleanupInterval := time.Duration(cfg.Print.CleanupInterval)
	go tasks.CleanTempTickets(ctx, cfg.Print.TempDir, cleanupInterval*time.Minute)

	return &PrintService{cfg: cfg}
}

func (service *PrintService) PrintToFileAndSend(content string) error {
	//maxFileSize := service.cfg.Print.MaxFileSize
	//allowedFormats := service.cfg.Print.AllowedFormats
	//tempDir := service.cfg.Print.TempDir
	printCommand := service.cfg.Print.Command
	//printTimeout := service.cfg.Print.Timeout
	//queueSize := service.cfg.Print.QueueSize

	ticketStorage, er := ticketfiles.NewTicketStorage(&service.cfg.Print)
	if er != nil {
		log.Fatalf("erreur lors de l'initialisation du stockage des tickets\n")
	}

	filePath, err := ticketStorage.CreateTempFile(content)
	if err != nil {
		log.Fatalf("erreur lors de la création du fichier temporaire : %v\n", err)
	}

	//ticket data pour tester
	ticketMetaData := ticketfiles.CreateTicketMetadata("ORD-100", "client", content, "XPRINTER D-200N", time.Now())

	logger.Log.Infof("\nFichier temporaire créé avec succès : %s\n", filePath)

	// TODO : validation
	//fileInfo, err := file.Stat()
	//if err != nil {
	//	_ = file.Close()
	//	return fmt.Errorf("erreur obtention info fichier : %v", err)
	//}
	//
	//err = file.Close()
	//if err != nil {
	//	return fmt.Errorf("erreur fermeture fichier :%v", err)
	//}
	//
	//if fileInfo.Size() > maxFileSize {
	//	return fmt.Errorf("erreur : taille fichier dépasse la limite autorisée de %d octets", maxFileSize)
	//}

	cmd := fmt.Sprintf(printCommand, filePath)
	output, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur exécution commande d'impression : %v, sortie : %s", err, string(output))
	}

	logger.Log.Infof("commande d'impression exécutée avec succès")
	ticketMetaData.PrintedAt = time.Now()
	err = ticketStorage.ArchiveTicket(ticketMetaData)
	if err != nil {
		return fmt.Errorf("erreur archivage ticket : %v", err)
	}

	//err = ticketStorage.CleanupTempFile(filePath)
	//if err != nil {
	//	return fmt.Errorf("erreur nettoyage fichier temporaire : %v", err)
	//}
	return nil

}

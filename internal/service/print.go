package service

import (
	"fmt"
	"os"
	"os/exec"
	"print-service/internal/config"
	"print-service/internal/logger"
)

type PrintService struct {
	cfg *config.Config
}

func NewPrintService(cfg *config.Config) *PrintService {
	return &PrintService{cfg: cfg}
}

func (service *PrintService) PrintToFileAndSend(content string) error {
	maxFileSize := service.cfg.Print.MaxFileSize
	//allowedFormats := service.cfg.Print.AllowedFormats
	tempDir := service.cfg.Print.TempDir
	printCommand := service.cfg.Print.Command
	//printTimeout := service.cfg.Print.Timeout
	//queueSize := service.cfg.Print.QueueSize

	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erreur création répertoire temporaire : %v", err)
	}

	filename := fmt.Sprintf("%s/print_job.txt", tempDir)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erreur création fichier : %v", err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("erreur fichier : %v", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("erreur obtention info fichier : %v", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("erreur fermeture fichier :%v", err)
	}

	if fileInfo.Size() > maxFileSize {
		return fmt.Errorf("erreur : taille fichier dépasse la limite autorisée de %d octets", maxFileSize)
	}

	cmd := fmt.Sprintf(printCommand, file.Name())
	output, err := exec.Command("powershell", "-Command", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur exécution commande d'impression : %v, sortie : %s", err, string(output))
	}

	logger.Log.Infof("commande d'impression exécutée avec succès")
	return nil

}

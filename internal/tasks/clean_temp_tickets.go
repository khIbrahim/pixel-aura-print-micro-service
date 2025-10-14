package tasks

import (
	"context"
	"os"
	"path/filepath"
	"print-service/internal/logger"
	"time"
)

func CleanTempTickets(ctx context.Context, tempDir string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Log.Info("nettoyage le s")
			cleanupOldFiles(tempDir, interval*time.Minute)
		case <-ctx.Done():
			logger.Log.Info("task de nettoyage des tickets temporaires arrêtée")
			return
		}
	}
}

func cleanupOldFiles(tempDir string, interval time.Duration) {
	files, err := os.ReadDir(tempDir)
	if err != nil {
		logger.Log.Errorf("erreur lecture dossier temporaire : %v", err)
		return
	}

	cutoff := time.Now().Add(-interval)

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			logger.Log.Errorf("erreur obtention info fichier temporaire %s : %v", file.Name(), err)
			continue
		}

		if !info.IsDir() && info.ModTime().Before(cutoff) {
			filePath := filepath.Join(tempDir, file.Name())
			err = os.Remove(filePath)
			if err != nil {
				logger.Log.Errorf("erreur suppression fichier temporaire %s : %v", filePath, err)
				continue
			}

			logger.Log.Infof("fichier temporaire supprimé : %s", filePath)
		} else if info.IsDir() {
			logger.Log.Infof("ignorer dossier dans temp : %s", file.Name())
		} else {
			logger.Log.Infof("Suppression du fichier temporaire dans %d secondes : %s", int(interval.Seconds()), file.Name())
		}
	}
}

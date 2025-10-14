package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"print-service/internal/config"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init(config *config.Config) error {
	level, err := logrus.ParseLevel(config.Logger.Level)
	if err != nil {
		return fmt.Errorf("level du log (%s) invalide : %v\n", config.Logger.Level, err)
	}

	outputFile := config.Logger.FilePath
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erreur lors de la création du dossier des logs : %v\n", err)
	}

	_, err = os.Stat(outputFile)
	if os.IsNotExist(err) {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du fichier des logs : %v\n", err)
		}

		defer file.Close()
		fmt.Printf("Le fichier des logs a bien été créé : %s\n", outputFile)
	} else if err != nil {
		return fmt.Errorf("erreur lors de la vérification du fichier des logs : %v\n", err)
	} else {
		fmt.Printf("Le fichier des logs détécté avec succès : %s\n", outputFile)
	}

	Log = logrus.New()
	Log.SetLevel(level)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		DisableColors: false,
	})
	Log.SetOutput(os.Stdout)
	Log.AddHook(&WriterHook{
		Writer: &lumberjack.Logger{
			Filename:   outputFile,
			MaxSize:    config.Logger.MaxSize,
			MaxBackups: config.Logger.MaxBackups,
			MaxAge:     config.Logger.MaxAge,
			Compress:   true,
		},
		LogLevels: logrus.AllLevels,
	})
	return nil
}

type WriterHook struct {
	Writer    *lumberjack.Logger
	LogLevels []logrus.Level
}

func (w WriterHook) Levels() []logrus.Level {
	return w.LogLevels
}

func (w WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	_, err = w.Writer.Write([]byte(line))
	return err
}

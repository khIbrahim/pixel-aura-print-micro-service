package config

import (
	"fmt"
	"runtime"
)

func validate(config *Config) error {
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("le port du serveur est invalide : %d", config.Server.Port)
	}

	// vu que je n'ai pas encore fait la gestion des os, ainsi que leur config etc on va juste valider que c'est un os avec le runtime
	if runtime.GOOS != "windows" {
		return fmt.Errorf("le système d'exploitation doit être windows pour le moment, détécté : %s", runtime.GOOS)
	}

	// jlaisse en TODO, jai pas d'inspi

	return nil
}

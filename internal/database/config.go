package database

import (
	"path/filepath"

	"github.com/Cufee/valorant-account-tracker-go/internal/config"
)

var localDatabasePath = filepath.Join(config.AppDataPath, "database", "main.db")

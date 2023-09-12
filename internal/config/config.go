package config

import (
	"os"
	"path/filepath"
)

var AppDataPath = filepath.Join(os.Getenv("LocalAppData"), "byvko-dev", "valorant-account-tracker")

package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Cufee/valorant-account-tracker-go/internal/config"
)

var localDatabasePath = filepath.Join(config.AppDataPath, "database", "main.db")

func init() {
	err := os.MkdirAll(filepath.Dir(localDatabasePath), 0600)
	if err != nil {
		log.Panicf("Failed to create database directory: %s", err)
	}
}

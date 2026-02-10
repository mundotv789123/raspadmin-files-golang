package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	AbsRootDir     string
	AllowedOrigins []string
)

func Init() {
	loadRootDir()
	loadCors()
}

func loadRootDir() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("Error loading .env file: ", err)
	}

	rootDir := os.Getenv("FILES_PATH")
	if rootDir == "" {
		rootDir = "./files"
	}
	AbsRootDir, _ = filepath.Abs(rootDir)
}

func loadCors() {
	originsString := os.Getenv("ALLOWED_ORIGINS")
	if originsString == "" {
		originsString = "http://localhost:3000"
	}
	AllowedOrigins = strings.Split(originsString, ",")
}

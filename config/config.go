package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	RootDir    string
	AbsRootDir string

	OriginsString  string
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

	RootDir = os.Getenv("FILES_PATH")
	if RootDir == "" {
		RootDir = "./files"
	}
	AbsRootDir, _ = filepath.Abs(RootDir)
}

func loadCors() {
	OriginsString = os.Getenv("ALLOWED_ORIGINS")
	if OriginsString != "" {
		OriginsString = "http://localhost:3000"
	}
	AllowedOrigins = strings.Split(OriginsString, ",")
}

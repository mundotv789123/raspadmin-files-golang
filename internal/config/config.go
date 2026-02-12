package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	AbsRootDir     string
	AllowedOrigins []string

	JwtSecret string

	AuthEnabled bool
	AppUsername string
	AppPassword string
)

func Init() {
	loadRootDir()
	loadCors()
	loadAuth()
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

func loadAuth() {
	authEnabledEnv := os.Getenv("AUTH_ENABLED")
	AuthEnabled = authEnabledEnv == "true"

	if !AuthEnabled {
		return
	}

	JwtSecret = os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		JwtSecret = uuid.New().String()
	}

	AppUsername = os.Getenv("APP_USERNAME")
	if AppUsername == "" {
		AppUsername = "admin"
	}

	AppPassword = os.Getenv("APP_PASSWORD")
	if AppPassword == "" {
		AppPassword = "admin"
	}
}

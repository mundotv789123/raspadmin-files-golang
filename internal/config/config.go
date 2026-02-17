package config

import (
	"log/slog"
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

	CacheDir    string
	CacheDirAds string

	DbFile string
)

func Init() {
	loadRootDir()
	loadCors()
	loadAuth()
	loadCache()
	loadDatabase()
}

func loadRootDir() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("Error loading .env file: ", err)
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

func loadCache() {
	CacheDir = os.Getenv("CACHE_DIR")
	if CacheDir == "" {
		CacheDir = "_cache"
	}
	CacheDirAds, _ = filepath.Abs(filepath.Join(AbsRootDir, CacheDir))
}

func loadDatabase() {
	DbFile = os.Getenv("DB_FILE")
	if DbFile == "" {
		DbFile = "database.db"
	}
}

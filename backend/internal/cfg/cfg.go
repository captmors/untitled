package cfg

import (
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

var (
	RootDir     string
	DatabaseUrl string
	JwtSecret   string

	// default
	LogDir         string
	LogDefaultFile string
	LogToFile      bool
	TestLogToFile  bool
	UploadDir      string
)

func init() {
	// init RootDir
	_, b, _, _ := runtime.Caller(0)
	RootDir = filepath.Join(filepath.Dir(filepath.Dir(b)), "..")

	// parse .env if exists locally
	envPath := filepath.Join(RootDir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	} else {
		log.Println(".env file not found, using system environment variables.")
	}

	DatabaseUrl = os.Getenv("DATABASE_URL")
	JwtSecret = os.Getenv("JWT_SECRET")

	// default
	LogDir = filepath.Join(RootDir, getEnvOr("LOG_DIR", "logs"))
	LogDefaultFile = getEnvOr("LOG_DEFAULT_FILE", "log.log")
	LogToFile = getEnvOrBool("LOG_TO_FILE", false)
	TestLogToFile = getEnvOrBool("TEST_LOG_TO_FILE", false)

	UploadDir = filepath.Join(RootDir, getEnvOr("UPLOAD_DIR", "uploads"))
}

func getEnvOr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvOrBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true" || value == "1" || value == "True" || value == "TRUE"
	}
	return fallback
}

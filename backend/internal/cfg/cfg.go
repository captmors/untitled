package cfg

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

var (
	RootDir     string
	JwtSecret   string

	// PostgreSQL
    PostgresUser     string
    PostgresPassword string
    PostgresDB       string

    // MongoDB
    MongoUser     string
    MongoPassword string

    // ElasticSearch
    ElasticUser     string
    ElasticPassword string

	// default
	LogDir            string
	LogDefaultFile    string
	LogToFile         bool
	TestLogToFile     bool
	UploadDir         string
	MaxUploadFileSize int64
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
		log.Warning(".env file not found, using system environment variables.")
	}

	JwtSecret = os.Getenv("JWT_SECRET")

	PostgresUser = os.Getenv("POSTGRES_USER")
    PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
    PostgresDB = os.Getenv("POSTGRES_DB")

    MongoUser = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
    MongoPassword = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

    ElasticUser = os.Getenv("ELASTICSEARCH_USERNAME")
    ElasticPassword = os.Getenv("ELASTICSEARCH_PASSWORD")

	// default
	LogDir = filepath.Join(RootDir, getEnvOr("LOG_DIR", "logs"))
	LogDefaultFile = getEnvOr("LOG_DEFAULT_FILE", "log.log")
	LogToFile = getEnvOrBool("LOG_TO_FILE", false)
	TestLogToFile = getEnvOrBool("TEST_LOG_TO_FILE", false)

	UploadDir = filepath.Join(RootDir, getEnvOr("UPLOAD_DIR", "uploads"))
	MaxUploadFileSize = getEnvOrInt64("MAX_UPLOAD_FILE_SIZE", 50 * 1024 * 1024)
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

func getEnvOrInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return fallback
}
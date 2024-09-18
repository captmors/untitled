// cfg/config.go

package cfg

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	RootDir     string
	DatabaseUrl string
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
}

// cfg/config.go

package cfg

import (
    "path/filepath"
    "runtime"
    "log"
    "os"
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

	// parse .env
    envPath := filepath.Join(RootDir, ".env")
    if err := godotenv.Load(envPath); err != nil {
        log.Fatal("Error loading .env file:", err)
    }

    DatabaseUrl = os.Getenv("DATABASE_URL")
}

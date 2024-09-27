package cfg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"untitled/internal/interfaces"
	"untitled/internal/users"
	"untitled/internal/users/mdl"
	tus "untitled/internal/utils/tusuploader"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	initLogging()
}

func Init() *gin.Engine {
	r := gin.Default()

	db := initDB()
	initApps(r, db)

	tus.InitTusUploader(r, tus.TusHandlerCfg{
		MaxFileSize: MaxUploadFileSize,
		UploadDir:   UploadDir,
	})

	return r
}

var (
	apps []interfaces.App
)

func initDB() *gorm.DB {
	DB, err := gorm.Open(postgres.Open(DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err := DB.AutoMigrate(&mdl.User{}); err != nil {
		log.Fatal("Error migrating database:", err)
	}

	return DB
}

func initApps(r *gin.Engine, db *gorm.DB) {
	usersApp := users.NewUserApp(db, []byte(JwtSecret))
	apps = append(apps, usersApp)

	for _, app := range apps {
		app.Init(r)
	}
}

// ENV:
// - LOG_TO_FILE: bool
func initLogging() *os.File {
	logDir := filepath.Join(LogDir)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	var file *os.File
	var err error

	if LogToFile {
		logFile := filepath.Join(logDir, LogDefaultFile)
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		multiWriter := io.MultiWriter(file, os.Stdout)
		log.SetOutput(multiWriter)
	} else {
		log.SetOutput(os.Stdout)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line), ""
		},
	})
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)

	return file
}

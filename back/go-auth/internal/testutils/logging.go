package testutils

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"untitled/internal/cfg"

	log "github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

// test case scope logging
// ENV:
// - TestLogToFile: bool
func InitTestLogging(testName string) (file *os.File) {
	logDir := filepath.Join(cfg.LogDir, "test")
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	if cfg.TestLogToFile {
		logFile := filepath.Join(logDir, testName+".log")
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("Failed to open test log file: %v", err)
		}
		multiWriter := io.MultiWriter(os.Stdout, file)
		log.SetOutput(multiWriter)
	} else {
		log.SetOutput(os.Stdout)
	}
	log.SetFormatter(&CustomFormatter{})
	log.SetLevel(log.InfoLevel)

	// header

	log.Println("========================================")

	_, fileName, line, ok := runtime.Caller(1)
	t := time.Now()
	timeFmt := "15:04:05 | 2006 02 01"
	if ok {
		log.Printf("| Source: %s:%d", fileName, line)
		log.Println("|---------------------------------------")
	}

	log.Printf("| Time:   %s", t.Format(timeFmt))
	log.Println("========================================")
	log.Println("")

	return
}

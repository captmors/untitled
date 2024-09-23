package main

import (
	"untitled/internal/cfg"

	log "github.com/sirupsen/logrus"
)

func main() {
	r := cfg.Init()
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

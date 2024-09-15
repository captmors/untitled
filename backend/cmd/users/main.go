package main

import (
	"log"
	"untitled/internal/cfg"
)

func main() {
	r := cfg.Init()
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

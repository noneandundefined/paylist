package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"paylist.server/config"
)

var (
	Version = "1.0.x"
)

func (s *httpServer) httpStart() error {
	port, err := strconv.Atoi(config.HttpServerPort)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	/* Initial HTTPx routes */
	routes := s.routes()

	fmt.Printf("\n[%v] [INFO] HTTP server started :%d\n", time.Now().Format("2006-01-02 15:04:05"), port)
	fmt.Printf("[%v] [INFO] Proccess PID: %d, Version: %s\n", time.Now().Format("2006-01-02 15:04:05"), os.Getpid(), Version)
	fmt.Printf("[%v] [INFO] Golang version: %s\n\n", time.Now().Format("2006-01-02 15:04:05"), runtime.Version())

	httpServe := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  90 * time.Second,
	}

	return httpServe.ListenAndServe()
}

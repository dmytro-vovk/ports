package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/dmytro-vovk/ports/services/rest-api-server/internal/boot"
)

func main() {
	configFile := flag.String("config", "rest-api-config.json", "Path to config file")

	flag.Parse()

	container, shutdown, err := boot.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := container.WebServer().Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Error running web server: %s", err)
	}

	shutdown()
}

package main

import (
	"flag"
	"log"

	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/boot"
)

func main() {
	configFile := flag.String("config", "storage-api-config.json", "Path to config file")

	flag.Parse()

	container, shutdown, err := boot.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := container.GRPCServer().Run(); err != nil {
		log.Printf("Error running GRPC server: %s", err)
	}

	shutdown()
}

package boot

import (
	"log"

	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/api"
	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/grpcserver"
	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/storage/memory"
)

func (c *Container) API() *api.API {
	const id = "API"

	if s, ok := c.get(id).(*api.API); ok {
		return s
	}

	s := api.New(c.Storage())

	log.Printf("Initialised %s", id)

	return c.set(id, s, nil).get(id).(*api.API)
}

func (c *Container) Storage() *memory.Storage {
	const id = "Storage"

	if s, ok := c.get(id).(*memory.Storage); ok {
		return s
	}

	s := memory.New()

	log.Printf("Initialised %s", id)

	return c.set(id, s, nil).get(id).(*memory.Storage)
}

func (c *Container) GRPCServer() *grpcserver.Server {
	const id = "GRPC Server"

	if s, ok := c.get(id).(*grpcserver.Server); ok {
		return s
	}

	s := grpcserver.New(c.Config().GRPCServer.Listen)

	log.Printf("Initialised %s", id)

	return c.set(id, s, s.Shutdown).get(id).(*grpcserver.Server)
}

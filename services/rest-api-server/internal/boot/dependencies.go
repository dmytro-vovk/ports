package boot

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dmytro-vovk/ports/services/rest-api-server/internal/api"
	"github.com/dmytro-vovk/ports/services/rest-api-server/internal/storage"
	"github.com/dmytro-vovk/ports/services/rest-api-server/internal/webserver"
	"github.com/dmytro-vovk/ports/services/rest-api-server/internal/webserver/router"
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

func (c *Container) Router() *router.Router {
	const id = "Router"

	if s, ok := c.get(id).(*router.Router); ok {
		return s
	}

	s := router.New(
		router.Route{
			Method:  http.MethodPut,
			Path:    "/port",
			Handler: c.API().UploadPorts,
		},
		router.Route{
			Method:  http.MethodGet,
			Path:    "/port/{id:[A-Z]+}",
			Handler: c.API().GetPortByID,
		},
	)

	log.Printf("Initialised %s", id)

	return c.set(id, s, nil).get(id).(*router.Router)
}

func (c *Container) Storage() *storage.Client {
	const id = "Storage"

	if s, ok := c.get(id).(*storage.Client); ok {
		return s
	}

	s := storage.New(c.Config().StorageAPI.Address)

	log.Printf("Initialised %s", id)

	return c.set(id, s, func() {
		if err := s.Close(); err != nil {
			log.Printf("Error shutting down %s: %s", id, err)
		}
	}).get(id).(*storage.Client)
}

func (c *Container) WebServer() *webserver.Server {
	const id = "Web Server"

	if s, ok := c.get(id).(*webserver.Server); ok {
		return s
	}

	s := webserver.New(c.Config().WebServer.Listen, c.Router())

	log.Printf("Initialised %s", id)

	return c.set(id, s, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		if err := s.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down %s: %s", id, err)
		}

		cancel()
	}).get(id).(*webserver.Server)
}

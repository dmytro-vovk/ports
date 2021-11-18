package webserver

import (
	"context"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	address string
	handler http.Handler
	server  *http.Server
	once    sync.Once
}

func New(address string, handler http.Handler) *Server {
	return &Server{
		address: address,
		handler: handler,
	}
}

func (s *Server) Run() (err error) {
	s.once.Do(func() {
		s.server = &http.Server{
			Addr:    s.address,
			Handler: s.handler,
		}

		log.Printf("Server starting on %s...", s.address)

		err = s.server.ListenAndServe()
	})

	return
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

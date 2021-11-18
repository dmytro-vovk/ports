package grpcserver

import (
	"log"
	"net"

	"github.com/dmytro-vovk/ports/services/protocol"
	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/api"
	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/storage/memory"
	"google.golang.org/grpc"
)

type Server struct {
	address string
	gs      *grpc.Server
}

func New(address string) *Server {
	return &Server{address: address}
}

func (s *Server) Run() error {
	conn, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	s.gs = grpc.NewServer()
	protocol.RegisterStorageServer(s.gs, api.New(memory.New()))

	log.Printf("Server starting on %s...", s.address)

	return s.gs.Serve(conn)
}

func (s *Server) Shutdown() { s.gs.GracefulStop() }

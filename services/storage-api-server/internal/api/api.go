package api

import (
	"context"

	"github.com/dmytro-vovk/ports/services/protocol"
	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/storage"
)

// API is a set of business logic handlers for our application
type API struct {
	protocol.UnimplementedStorageServer
	db storage.Storage
}

func New(db storage.Storage) *API {
	return &API{
		db: db,
	}
}

func (a *API) Store(server protocol.Storage_StoreServer) error {
	for {
		request, err := server.Recv()
		if err != nil {
			return err
		}

		ctx := server.Context()

		for key, value := range request.List {
			if err := a.db.Insert(ctx, key, value); err != nil {
				return err
			}
		}
	}
}

func (a *API) Get(ctx context.Context, port *protocol.GetPortRequest) (*protocol.Data, error) {
	return a.db.Get(ctx, port.Name)
}

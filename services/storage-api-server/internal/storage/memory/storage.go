package memory

import (
	"context"
	"database/sql"
	"sync"

	"github.com/dmytro-vovk/ports/services/protocol"
	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/storage"
)

// Storage is in-memory implementation of storage.Storage interface for testing purposes
type Storage struct {
	items  map[string]*protocol.Data
	itemsM sync.Mutex
}

var _ storage.Storage = &Storage{}

func New() *Storage {
	return &Storage{
		items: make(map[string]*protocol.Data),
	}
}

func (d *Storage) Insert(_ context.Context, key string, value *protocol.Data) error {
	d.itemsM.Lock()
	defer d.itemsM.Unlock()

	d.items[key] = value

	return nil
}

func (d *Storage) Get(_ context.Context, key string) (*protocol.Data, error) {
	d.itemsM.Lock()
	defer d.itemsM.Unlock()

	item, ok := d.items[key]
	if ok {
		return item, nil
	}

	return nil, sql.ErrNoRows // Using sql.ErrNoRows, as real storage will most probably use it
}

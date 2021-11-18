package storage

import (
	"context"

	"github.com/dmytro-vovk/ports/services/protocol"
)

type Storage interface {
	Insert(ctx context.Context, key string, value *protocol.Data) error
	Get(ctx context.Context, key string) (*protocol.Data, error)
}

package memory_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dmytro-vovk/ports/services/protocol"
	"github.com/dmytro-vovk/ports/services/storage-api-server/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	s := memory.New()
	require.NoError(t, s.Insert(context.Background(), "key1", &protocol.Data{Name: "A"}))

	data, err := s.Get(context.Background(), "key1")
	require.NoError(t, err)

	assert.Equal(t, &protocol.Data{Name: "A"}, data)

	_, err = s.Get(context.Background(), "non existing key")
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

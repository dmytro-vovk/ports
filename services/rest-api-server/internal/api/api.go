package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dmytro-vovk/ports/services/protocol"
	"github.com/dmytro-vovk/ports/services/rest-api-server/internal/stream"
	"github.com/gorilla/mux"
)

// API is a set of business logic handlers for our application
type API struct {
	storage StorageClient
}

const bufferSize = 10

func New(client StorageClient) *API {
	return &API{
		storage: client,
	}
}

func (a *API) UploadPorts(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.Closer) {
		_ = Body.Close() //nolint:staticcheck // No op
	}(r.Body)
	// Get stream of portData from request body
	dataStream, err := stream.Scan(r.Body, portData{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	// Buffer for accumulating items
	buffer := map[string]*protocol.Data{}
	// Read the stream
	for item := range dataStream {
		if item.Error != nil {
			http.Error(w, item.Error.Error(), http.StatusBadRequest)

			return
		}

		buffer[item.Key] = item.Value.(*portData).AsProtocolData()

		if len(buffer) < bufferSize {
			continue
		}
		// Send buffered data
		if err := a.storage.Store(context.Background(), &protocol.StorePortRequest{List: buffer}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
		// Clear the buffer
		buffer = map[string]*protocol.Data{}
	}
	// Send the remaining data
	if len(buffer) != 0 {
		if err := a.storage.Store(context.Background(), &protocol.StorePortRequest{List: buffer}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}
}

func (a *API) GetPortByID(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		return
	}

	data, err := a.storage.Get(context.Background(), &protocol.GetPortRequest{Name: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

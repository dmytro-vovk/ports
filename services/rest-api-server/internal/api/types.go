package api

import (
	"context"

	"github.com/dmytro-vovk/ports/services/protocol"
)

type StorageClient interface {
	Store(ctx context.Context, request *protocol.StorePortRequest) error
	Get(ctx context.Context, request *protocol.GetPortRequest) (*protocol.Data, error)
}

type portData struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
	Coordinates []float32 `json:"coordinates"`
}

func (d portData) AsProtocolData() *protocol.Data {
	data := protocol.Data{
		Name:     d.Name,
		City:     d.City,
		Province: d.Province,
		Country:  d.Country,
		Alias:    d.Alias,
		Regions:  d.Regions,
		Timezone: d.Timezone,
		Unlocs:   d.Unlocs,
		Code:     d.Code,
	}

	if len(d.Coordinates) == 2 {
		data.Lat, data.Lon = d.Coordinates[0], d.Coordinates[1]
	}

	return &data
}

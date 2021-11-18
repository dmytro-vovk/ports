//go:build integration

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/dmytro-vovk/ports/services/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRESTServer(t *testing.T) {
	feedData(t)

	resp, err := http.Get("http://rest-api:5001/port/AIRBY")
	require.NoError(t, err)

	var body protocol.Data

	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

	assert.Equal(t, protocol.Data{
		Name:     "Road Bay",
		City:     "Road Bay",
		Province: "",
		Country:  "Anguilla",
		Alias:    nil,
		Regions:  nil,
		Timezone: "America/Anguilla",
		Unlocs:   []string{"AIRBY"},
		Code:     "24821",
		Lat:      -63.065327,
		Lon:      18.220749,
	}, body)
}

func feedData(t *testing.T) {
	f, err := os.Open("internal/stream/test-data/ports.json")
	require.NoError(t, err)

	defer f.Close()

	req, err := http.NewRequest(http.MethodPut, "http://rest-api:5001/port", f)
	require.NoError(t, err)

	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
}

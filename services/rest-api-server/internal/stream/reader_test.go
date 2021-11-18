package stream

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// FIXME happy path only
func TestScanner(t *testing.T) {
	type record struct {
		Name        string        `json:"name"`
		City        string        `json:"city"`
		Province    string        `json:"province"`
		Country     string        `json:"country"`
		Alias       []string      `json:"alias"`
		Regions     []interface{} `json:"regions"` // Array of unknown type
		Coordinates []float32     `json:"coordinates"`
		Timezone    string        `json:"timezone"`
		Unlocs      []string      `json:"unlocs"`
		Code        string        `json:"code"`
	}

	const testData = "test-data/ports.json"

	// Read reference data
	var records map[string]*record

	data, err := ioutil.ReadFile(testData)
	require.NoError(t, err)

	require.NoError(t, json.Unmarshal(data, &records))

	f, err := os.Open(testData)
	require.NoError(t, err)
	defer f.Close()

	c, err := Scan(f, &record{})
	require.NoError(t, err)

	var n int // Records counter

	for r := range c {
		require.NoError(t, r.Error, "Unexpected error reading records stream")
		require.Equal(t, records[r.Key], r.Value, "A record not found in reference data")
		n++
	}

	assert.Equal(t, len(records), n, "Wrong number of records")
}

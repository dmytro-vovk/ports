package stream

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestScannerHappyPath(t *testing.T) {
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

func TestScannerErrors(t *testing.T) {
	testCases := []struct {
		name    string
		source  io.Reader
		initErr error
		scanErr string
	}{
		{
			name:    "Empty source",
			source:  bytes.NewBuffer([]byte{}),
			initErr: io.EOF,
		},
		{
			name:    "Empty object",
			source:  bytes.NewBuffer([]byte(`{}`)),
			initErr: nil,
		},
		{
			name:    "Not an object",
			source:  bytes.NewBuffer([]byte(`[]`)),
			initErr: errors.New("expected '{', got json.Delim([)"),
		},
		{
			name:    "Invalid key",
			source:  bytes.NewBuffer([]byte(`{0:`)),
			scanErr: "invalid character '0'",
		},
		{
			name:    "Invalid value",
			source:  bytes.NewBuffer([]byte(`{"a":`)),
			scanErr: "EOF",
		},
		{
			name:    "Invalid value type",
			source:  bytes.NewBuffer([]byte(`{"a": false}`)),
			scanErr: "json: cannot unmarshal bool into Go value of type stream.record",
		},
	}

	for i := range testCases {
		tt := testCases[i]

		t.Run(tt.name, func(t *testing.T) {
			c, err := Scan(tt.source, &record{})
			assert.Equal(t, tt.initErr, err)

			if err == nil {
				for val := range c {
					if tt.scanErr == "" {
						assert.NoError(t, val.Error)
					} else if assert.Error(t, val.Error) {
						assert.Equal(t, tt.scanErr, val.Error.Error())
					}
				}
			}
		})
	}
}

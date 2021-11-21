package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// Record represents the result of scanning JSON object key/value pair
type Record struct {
	Key   string
	Value interface{}
	Error error
}

// Scan tries to read records from the reader asynchronously sending them via channel
func Scan(r io.Reader, data interface{}) (<-chan Record, error) {
	// Get data type, literal and pointer are both acceptable
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}
	// Create decoder
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields() // Stop on an unexpected field
	// Read object opening brace
	openingBrace, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	// If it is not an object, return an error
	if openingBrace != json.Delim('{') {
		return nil, fmt.Errorf("expected '{', got %[1]T(%[1]v)", openingBrace)
	}
	// Results channel
	resultsChan := make(chan Record)
	// Start decoding
	go func(c chan Record, d *json.Decoder) {
		defer close(c)
		// Loop over the object properties
		for d.More() {
			// Read object key
			token, err := d.Token()
			if err != nil {
				c <- Record{Error: err}

				break
			}
			// Create new instance of data type
			instance := reflect.New(dataType).Interface()
			// Decode the value
			if err = d.Decode(&instance); err != nil {
				c <- Record{Error: err}
				// Stop on error
				break
			}
			// Send it
			c <- Record{
				Key:   token.(string),
				Value: instance,
			}
		}
	}(resultsChan, decoder)

	return resultsChan, nil
}

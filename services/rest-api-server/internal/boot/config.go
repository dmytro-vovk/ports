package boot

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type Config struct {
	WebServer struct {
		Listen string `json:"listen"`
	} `json:"web_server"`
	StorageAPI struct {
		Address string `json:"address"`
	} `json:"storage_api"`
}

func readConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "reading %s", fileName)
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrapf(err, "decoding %s", fileName)
	}

	return &cfg, nil
}

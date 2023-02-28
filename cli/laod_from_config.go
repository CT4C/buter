package cli

import (
	"encoding/json"
	"os"
)

func loadJSONConfig(filename string) ([]UserConfig, error) {
	c := make([]UserConfig, 0)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	s, _ := file.Stat()

	data := make([]byte, s.Size())

	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &c)

	return c, err
}

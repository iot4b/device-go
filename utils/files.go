package utils

import (
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"io"
	"os"
)

func SaveFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	log.Debug(string(data))
	return nil
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open(%s)", path)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrapf(err, "io.ReadAll(%s)", path)
	}
	return bytes, nil
}

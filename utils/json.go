package utils

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var (
	ErrUnmarshal = errors.New("file unmarshal error")
)

func ReadJSONFile(path string, to any) error {
	fileData, err := ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "ReadFile(%s)", path)
	}
	err = json.Unmarshal(fileData, to)
	if err != nil {
		return errors.Wrapf(ErrUnmarshal, "json.Unmarshal(%s, to): %s", path, err.Error())
	}
	return nil
}

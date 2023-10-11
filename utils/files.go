package utils

import (
	"bufio"
	"encoding/json"
	"github.com/markgenuine/ever-client-go/domain"
	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
	"io"
	"os"
)

var (
	ErrUnmarshal = errors.New("file unmarshal error")
)

func JsonMapToStruct(input interface{}, out interface{}) error {
	d, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, &out)
}

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

func ReadKeysFile(path string) (keys domain.KeyPair, err error) {
	err = ReadJSONFile(path, &keys)
	return
}

func WriteToStdout(data []byte) error {
	// иначе выводим в Stdout
	writer := bufio.NewWriter(os.Stdout)
	if _, err := writer.WriteString(string(data)); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func OutJson(data interface{}) error {
	// формируем json на выход
	result, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "json.Marshal(data)")
	}
	writer := bufio.NewWriter(os.Stdout)
	if _, err := writer.WriteString(string(result)); err != nil {
		return errors.Wrap(err, "WriteString(out)")
	}
	writer.Flush()
	return nil
}

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

package utils

import (
	"encoding/json"
	"io"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/ndmsystems/golog"
	"github.com/pkg/errors"
)

const filesDir = "iot4bd"

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

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetFilesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		u, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		home = u.HomeDir
	}
	dir := filepath.Join(home, ".config", filesDir)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
	return dir
}

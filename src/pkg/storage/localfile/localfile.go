package localfile

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"hms/gateway/pkg/errors"
	"log"
	"os"
)

type Config struct {
	BasePath string
	Depth    uint8
}

type LocalFileStorage struct {
	basePath string
	depth    uint8
}

func Init(config *Config) (*LocalFileStorage, error) {
	if len(config.BasePath) == 0 {
		return nil, fmt.Errorf("BasePath is empty")
	}

	if config.Depth == 0 {
		config.Depth = 1
	}

	if config.BasePath[len(config.BasePath)-1] != '/' {
		config.BasePath += "/"
	}

	_, err := os.Stat(config.BasePath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(config.BasePath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &LocalFileStorage{
		basePath: config.BasePath,
		depth:    config.Depth,
	}, nil
}

func (s *LocalFileStorage) Add(data []byte) (id *[32]byte, err error) {
	id = s.idByContent(&data)

	err = s.writeFile(id, &data)

	return
}

func (s *LocalFileStorage) idByContent(data *[]byte) *[32]byte {
	h := sha3.Sum256(*data)
	return &h
}

func (s *LocalFileStorage) ReplaceWithId(id *[32]byte, data []byte) (err error) {
	return s.AddWithId(id, data)
}

func (s *LocalFileStorage) AddWithId(id *[32]byte, data []byte) (err error) {
	err = s.writeFile(id, &data)
	return
}

func (s *LocalFileStorage) Get(id *[32]byte) (data []byte, err error) {
	idStr := hex.EncodeToString(id[:])

	path := s.filepath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return nil, errors.IsNotExist
	}
	data, err = os.ReadFile(path)

	return
}

func (s *LocalFileStorage) writeFile(id *[32]byte, data *[]byte) (err error) {
	idStr := hex.EncodeToString(id[:])

	path := s.dirpath(idStr)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	filepath := s.filepath(idStr)
	err = os.WriteFile(filepath, *data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *LocalFileStorage) dirpath(id string) (path string) {
	path = s.basePath
	i := 0
	for i < int(s.depth)*2 {
		path += id[i:i+2] + "/"
		i += 2
	}
	return path
}

func (s *LocalFileStorage) filepath(id string) (path string) {
	path = s.basePath
	i := 0
	for i < int(s.depth)*2 {
		path += id[i:i+2] + "/"
		i += 2
	}
	return path + id
}

func (s *LocalFileStorage) Clean() (err error) {
	if s.basePath == "/" {
		log.Panicln("Can not clean the base folder is root!")
	}

	_, err = os.Stat(s.basePath)
	if err != nil {
		return nil
	}

	if err = os.RemoveAll(s.basePath); err != nil {
		return err
	}

	return nil
}

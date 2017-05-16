package cache

import (
	"os"
	"path"
	"io/ioutil"
)

type Cache interface {
	Set(string, []byte) (string, error)
	Get(string) ([]byte, error)
	GetName(string) string
	Exists(string) bool
}

type cache struct {
	cacheDir	string
	nameConvert func(name string) string
}

func (c cache) convert(key string) string {
	if c.nameConvert != nil {
		return c.nameConvert(key)
	}
	return key
}

func (c cache) path(name string) string {
	return path.Join(c.cacheDir, c.convert(name))
}

func New(cacheDir string, conv func(name string) string) (Cache, error){
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	if conv != nil {
		return &cache{cacheDir:cacheDir, nameConvert: conv}, nil
	}
	return &cache{cacheDir:cacheDir}, nil
}

func (c cache) Exists(key string) bool {
	if _, err := os.Stat(c.path(key)); os.IsExist(err) {
		return true
	}
	return false
}

func (c cache) Set(key string, data []byte) (string, error) {
	return c.path(key), ioutil.WriteFile(c.path(key), data, 0644)
}

func (c cache) Get(key string) ([]byte, error) {
	return ioutil.ReadFile(c.path(key))
}

func (c cache) GetName(key string) string {
	return c.path(key)
}
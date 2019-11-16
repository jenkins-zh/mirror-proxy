package pkg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// CacheServer is the interface for saving a cache item
type CacheServer interface {
	Load(string) string
	Save(string, string)
}

// FileSystemCacheServer save the cache into a filesystem
type FileSystemCacheServer struct {
	FileName string

	cache map[string]string
}

// Load load the key from a file
func (c *FileSystemCacheServer) Load(key string) (val string) {
	if err := c.parse(); err == nil {
		val = c.cache[key]
	}
	return
}

// Save save the key into a file
func (c *FileSystemCacheServer) Save(key string, val string) (err error) {
	c.cache[key] = val
	var data []byte

	if data, err = yaml.Marshal(c.cache); err == nil {
		err = ioutil.WriteFile(c.FileName, data, 0644)
	}
	return
}

func (c *FileSystemCacheServer) parse() (err error) {
	var data []byte
	c.cache = make(map[string]string, 0)

	if data, err = ioutil.ReadFile(c.FileName); err == nil {
		err = yaml.Unmarshal(data, c.cache)
	}
	return nil
}


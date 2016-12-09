package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const session = ".consul_lock_id"

type Cache struct {
	basePath string
}

func New(directory string) (*Cache, error) {
	session := filepath.Join(directory, session)
	cache := &Cache{basePath: directory}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0777)
		sessionPath := filepath.Join(cache.basePath, session)
		if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
			return cache, cache.UpdateSession("")
		}
	}

	return cache, nil
}

func (c *Cache) Put(key string, value []byte) error {
	cacheDir := filepath.Join(c.basePath, key)
	return ioutil.WriteFile(cacheDir, value, 0644)
}

func (c *Cache) Get(key string) ([]byte, error) {
	cacheDir := filepath.Join(c.basePath, key)
	return ioutil.ReadFile(cacheDir)
}

// Get Session ID
func (c *Cache) GetSession() (string, error) {
	b, err := c.Get(session)
	return string(b), err
}

// Update Session ID
func (c *Cache) UpdateSession(id string) error {
	return c.Put(session, []byte(id))
}

package pgcache

import (
	"github.com/tanlosav/pg-cache/internal/configuration"
	"github.com/tanlosav/pg-cache/internal/db"
)

type Cache struct {
	config *configuration.Configuration
	driver *db.Driver
}

func NewCache(config *configuration.Configuration, driver *db.Driver) *Cache {
	return &Cache{
		config: config,
		driver: driver,
	}
}

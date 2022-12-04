package pgcache

import (
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type Cache struct {
	config configuration.Configuration
	driver *Driver
}

func NewCache(config configuration.Configuration) *Cache {
	driver := NewDriver(config)
	driver.Init()

	return &Cache{
		config: config,
		driver: driver,
	}
}

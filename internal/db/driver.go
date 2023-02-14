package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type Driver struct {
	config *configuration.Configuration
	Db     *sql.DB
}

func NewDriver(config *configuration.Configuration) *Driver {
	return &Driver{
		config: config,
	}
}

func (d *Driver) Connect() {
	addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", d.config.Db.User, d.config.Db.Password, d.config.Db.Host, d.config.Db.Name)
	db, err := sql.Open("postgres", addr)

	if err != nil {
		panic(err)
	}

	d.Db = db
}

// func (cache *Cache) Get(key string) (string, error) {
// 	var document string

// 	if err := cache.Db.QueryRow("select document from cache where key = $1", key).Scan(&document); err != nil {
// 		return "", err
// 	}

// 	return document, nil
// }

// func (cache *Cache) Create(key string, document []byte) error {
// 	_, err := cache.Db.Exec("insert into cache(key, document) values($1, $2::json)", key, document)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (cache *Cache) Update(key string, document []byte) error {
// 	_, err := cache.Db.Exec("insert into cache(key, document) values($1, $2::json) ON CONFLICT (key) DO UPDATE set document = $2::json", key, document)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (cache *Cache) Delete(key string) error {
// 	_, err := cache.Db.Exec("delete from cache where key = $1", key)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (cache *Cache) Clean() error {
// 	_, err := cache.Db.Exec("delete from cache")

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

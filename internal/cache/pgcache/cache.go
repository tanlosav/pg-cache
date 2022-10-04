package pgcache

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type Cache struct {
	User     string
	Password string
	Name     string
	Host     string
	Db       *sql.DB
}

func NewCache(user string, password string, name string, host string) *Cache {
	return &Cache{
		User:     user,
		Password: password,
		Name:     name,
		Host:     host,
	}
}

func (cache *Cache) Connect() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cache.User, cache.Password, cache.Host, cache.Name)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	cache.Db = db
}

func (cache *Cache) Get(key string) (string, error) {
	rows, err := cache.Db.Query("select document from cache where key = $1", key)

	if err != nil {
		return "", err
	}

	defer rows.Close()

	for rows.Next() {
		var document string
		if err := rows.Scan(&document); err != nil {
			return "", err
		}

		return document, nil
	}

	return "", errors.New("not found")
}

func (cache *Cache) Create(key string, document []byte) error {
	_, err := cache.Db.Exec("insert into cache(key, document) values($1, $2::json)", key, document)

	if err != nil {
		return err
	}

	return nil
}

func (cache *Cache) Update(key string, document []byte) error {
	_, err := cache.Db.Exec("insert into cache(key, document) values($1, $2::json) ON CONFLICT (key) DO UPDATE set document = $2::json", key, document)

	if err != nil {
		return err
	}

	return nil
}

func (cache *Cache) Delete(key string) error {
	_, err := cache.Db.Exec("delete from cache where key = $1", key)

	if err != nil {
		return err
	}

	return nil
}

func (cache *Cache) Clean() error {
	_, err := cache.Db.Exec("delete from cache")

	if err != nil {
		return err
	}

	return nil
}

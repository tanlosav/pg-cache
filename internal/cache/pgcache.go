package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type Cache struct {
	config     configuration.Configuration
	dbSettings map[string]configuration.Bucket
	db         *sql.DB
	tx         *sql.Tx
}

const (
	CACHE_SETTINGS_TABLE string = "pg_cache_settings"
	CACHE_LOCKS_SETTINGS int    = 1001
)

func NewPgCache(config configuration.Configuration) *Cache {
	return &Cache{
		config: config,
	}
}

func (cache *Cache) Connect() {
	addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cache.config.Db.User, cache.config.Db.Password, cache.config.Db.Host, cache.config.Db.Name)
	db, err := sql.Open("postgres", addr)

	if err != nil {
		panic(err)
	}

	cache.db = db
}

func (cache *Cache) Init() {
	cache.createDatabaseSchema()

	tx, err := cache.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	cache.tx = tx

	cache.lockSettingsTable()
	defer cache.tx.Rollback()

	cache.getDatabaseSettings()

	if len(cache.dbSettings) == 0 {
		log.Debug().Msg("Register cache settings to database.")
		cache.createPartitions()
		cache.storeDatabaseSettings()
		cache.tx.Commit()
	} else if !reflect.DeepEqual(cache.config.Cache.Buckets, cache.dbSettings) {
		panic("Please check cache settings. Database contains different configuration of the cache.")
	}
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

// createDatabaseSchema create database schema.
func (cache *Cache) createDatabaseSchema() {
	_, err := cache.db.Exec("create table if not exists " + CACHE_SETTINGS_TABLE + "(bucket varchar(128) not null primary key, settings varchar(1024))")
	if err != nil {
		panic(err)
	}
}

// getDatabaseSettings get stored settings for all bucket from database.
func (cache *Cache) getDatabaseSettings() {
	rows, err := cache.tx.Query("select bucket, settings from " + CACHE_SETTINGS_TABLE)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cache.dbSettings = make(map[string]configuration.Bucket)

	for rows.Next() {
		var bucket string
		var opts string

		err = rows.Scan(&bucket, &opts)
		if err != nil {
			panic(err)
		}

		var bucketOpts configuration.Bucket
		err = json.Unmarshal([]byte(opts), &bucketOpts)
		if err != nil {
			panic(err)
		}

		cache.dbSettings[bucket] = bucketOpts
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Printf("Database configuration: %+v", cache.dbSettings)
}

// storeDatabaseSettings store cache settings to database.
func (cache *Cache) storeDatabaseSettings() {
	for bucket, settings := range cache.config.Cache.Buckets {
		value, err := json.Marshal(settings)
		if err != nil {
			panic(err)
		}

		_, err = cache.tx.Exec("insert into "+CACHE_SETTINGS_TABLE+"(bucket, settings) values($1, $2)", bucket, value)
		if err != nil {
			panic(err)
		}
	}
}

// createPartitions create partitions for each bucket.
func (cache *Cache) createPartitions() {
	for bucket, settings := range cache.config.Cache.Buckets {
		for i := 0; i < settings.Sharding.PartitionsCount; i++ {
			cache.createPartition(bucket, settings.KeysCount, settings.Sharding.Partition, i)
		}
	}
}

// createPartition create single partition for the bucket.
func (cache *Cache) createPartition(bucket string, keysCount int, partition string, index int) {
	primaryKeyColumns := make([]string, 0, keysCount)
	partitionName := bucket + "_" + partition + strconv.Itoa(index)

	stmt := "create table " + partitionName + "("
	for i := 0; i < keysCount; i++ {
		primaryKeyColumn := "key_" + strconv.Itoa(i)
		primaryKeyColumns = append(primaryKeyColumns, primaryKeyColumn)
		stmt += primaryKeyColumn + " varchar not null,"
	}
	stmt += "document text not null,"
	stmt += "exp numeric not null,"
	stmt += "constraint " + partitionName + "_pk primary key (" + strings.Join(primaryKeyColumns, ",") + ")"
	stmt += ")"

	log.Debug().Msg("Create partition '" + partitionName + "' with statement: " + stmt)

	_, err := cache.db.Exec(stmt)
	if err != nil {
		panic(err)
	}
}

// lockSettingsTable acquire lock for settings table to check configuration and create partitions.
func (cache *Cache) lockSettingsTable() {
	log.Debug().Msg("Acquire lock on settings table.")

	_, err := cache.tx.Exec("select pg_advisory_xact_lock(" + strconv.Itoa(CACHE_LOCKS_SETTINGS) + ")")
	if err != nil {
		panic(err)
	}
}

package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type Schema struct {
	config *configuration.Configuration
	driver *Driver
}

const (
	CACHE_SETTINGS_TABLE string = "cache_settings"
	CACHE_LOCKS_SETTINGS int    = 1001
)

func NewSchema(config *configuration.Configuration, driver *Driver) *Schema {
	return &Schema{
		config: config,
		driver: driver,
	}
}

func (pm *Schema) Init() {
	pm.createDatabaseSchema()

	tx, err := pm.driver.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	pm.lockSettingsTable(tx)
	dbSettings := pm.getDatabaseSettings(tx)

	if len(dbSettings) == 0 {
		log.Debug().Msg("Register cache settings to database.")
		pm.createPartitions()
		pm.storeDatabaseSettings(tx)
		tx.Commit()
	} else if !reflect.DeepEqual(pm.config.Cache.Buckets, dbSettings) {
		panic("Please check cache settings. Database contains different configuration of the cache.")
	}
}

func (pm *Schema) createDatabaseSchema() {
	_, err := pm.driver.db.Exec("create table if not exists " + CACHE_SETTINGS_TABLE + "(bucket varchar(128) not null primary key, settings varchar(1024))")
	if err != nil {
		panic(err)
	}
}

func (pm *Schema) getDatabaseSettings(tx *sql.Tx) map[string]configuration.Bucket {
	rows, err := tx.Query("select bucket, settings from " + CACHE_SETTINGS_TABLE)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	dbSettings := make(map[string]configuration.Bucket)

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

		dbSettings[bucket] = bucketOpts
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Printf("Database configuration: %+v", dbSettings)

	return dbSettings
}

func (pm *Schema) storeDatabaseSettings(tx *sql.Tx) {
	for bucket, settings := range pm.config.Cache.Buckets {
		value, err := json.Marshal(settings)
		if err != nil {
			panic(err)
		}

		_, err = tx.Exec("insert into "+CACHE_SETTINGS_TABLE+"(bucket, settings) values($1, $2)", bucket, value)
		if err != nil {
			panic(err)
		}
	}
}

func (pm *Schema) createPartitions() {
	for bucket, settings := range pm.config.Cache.Buckets {
		for i := 0; i < settings.Sharding.PartitionsCount; i++ {
			pm.createPartition(bucket, settings.KeysCount, settings.Sharding.Partition, i)
		}
	}
}

func (pm *Schema) createPartition(bucket string, keysCount int, partition string, index int) {
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

	_, err := pm.driver.db.Exec(stmt)
	if err != nil {
		panic(err)
	}
}

func (pm *Schema) lockSettingsTable(tx *sql.Tx) {
	log.Debug().Msg("Acquire lock on settings table.")

	_, err := tx.Exec("select pg_advisory_xact_lock(" + strconv.Itoa(CACHE_LOCKS_SETTINGS) + ")")
	if err != nil {
		panic(err)
	}
}

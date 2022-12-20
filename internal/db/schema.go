package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

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

func (s *Schema) Init() {
	s.createDatabaseSchema()

	tx, err := s.driver.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	s.lockSettingsTable(tx)
	dbSettings := s.getDatabaseSettings(tx)

	if len(dbSettings) == 0 {
		log.Debug().Msg("Register cache settings to database.")
		s.createBuckets()
		s.storeDatabaseSettings(tx)
		tx.Commit()
	} else if !reflect.DeepEqual(s.config.Cache.Buckets, dbSettings) {
		panic("Please check cache settings. Database contains different configuration of the cache.")
	}
}

func (s *Schema) lockSettingsTable(tx *sql.Tx) {
	log.Debug().Msg("Acquire lock on settings table.")

	_, err := tx.Exec("select pg_advisory_xact_lock(" + strconv.Itoa(CACHE_LOCKS_SETTINGS) + ")")
	if err != nil {
		panic(err)
	}
}

func (s *Schema) createDatabaseSchema() {
	_, err := s.driver.db.Exec("create table if not exists " + CACHE_SETTINGS_TABLE + "(bucket varchar(128) not null primary key, settings varchar(1024))")
	if err != nil {
		panic(err)
	}
}

func (s *Schema) getDatabaseSettings(tx *sql.Tx) map[string]configuration.Bucket {
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

func (s *Schema) storeDatabaseSettings(tx *sql.Tx) {
	for bucket, settings := range s.config.Cache.Buckets {
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

func (s *Schema) createBuckets() {
	for bucket, settings := range s.config.Cache.Buckets {
		for partitionNumber := 0; partitionNumber < settings.Sharding.PartitionsCount; partitionNumber++ {
			s.createBucket(bucket, settings.KeysCount, partitionNumber, settings.Eviction)
		}
	}
}

func (s *Schema) createBucket(bucket string, keysCount int, partitionNumber int, eviction configuration.Eviction) {
	primaryKeyColumns := make([]string, 0, keysCount)
	tableName := bucket + "_" + strconv.Itoa(partitionNumber)

	stmt := "create table " + tableName + "("
	for i := 0; i < keysCount; i++ {
		primaryKeyColumn := "key_" + strconv.Itoa(i)
		primaryKeyColumns = append(primaryKeyColumns, primaryKeyColumn)
		stmt += primaryKeyColumn + " varchar not null,"
	}

	stmt += "document text not null,"
	stmt += "exp numeric not null,"
	stmt += "constraint " + tableName + "_pk primary key (" + strings.Join(primaryKeyColumns, ",") + ",exp)"
	stmt += ")"

	if configuration.EVICTION_POLICY_TRUNCATE == eviction.Policy {
		stmt += " PARTITION BY RANGE (exp)"
	}

	log.Debug().Msg("Create table '" + tableName + "' with statement: " + stmt)

	_, err := s.driver.db.Exec(stmt)
	if err != nil {
		panic(err)
	}

	if configuration.EVICTION_POLICY_TRUNCATE == eviction.Policy {
		s.createTablePartitions(tableName, eviction.PartitionTimeRange, eviction.ActualPartitionsCount)
	}
}

func (s *Schema) createTablePartitions(tableName string, timeRange int, count int) {
	now := time.Now().Unix()
	timeRangeInt64 := int64(timeRange)
	countInt64 := int64(count)

	for i := int64(0); i < countInt64; i++ {
		startValue, endValue := getPartitionBorders(now, timeRangeInt64, i)

		from := strconv.FormatInt(startValue, 10)
		to := strconv.FormatInt(endValue, 10)
		partitionName := tableName + "_" + from + "_" + to

		stmt := "create table " + partitionName + " partition of " + tableName + " for values from (" + from + ") to (" + to + ")"
		log.Debug().Msg("Create partition '" + partitionName + "' with statement: " + stmt)

		_, err := s.driver.db.Exec(stmt)
		if err != nil {
			panic(err)
		}
	}
}

func getPartitionBorders(now int64, timeRange int64, offset int64) (int64, int64) {
	start := now/timeRange*timeRange + timeRange*offset
	end := start + timeRange

	return start, end
}

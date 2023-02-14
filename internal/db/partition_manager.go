package db

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/package/set"
)

type PartitionManager struct {
	driver *Driver
}

type partitionBorders struct {
	startValue int64
	endValue   int64
}

func NewPartitionManager(driver *Driver) *PartitionManager {
	return &PartitionManager{
		driver: driver,
	}
}

func (p *PartitionManager) createNew(tableName string, bordersList []partitionBorders) {
	for _, borders := range bordersList {
		from := strconv.FormatInt(borders.startValue, 10)
		to := strconv.FormatInt(borders.endValue, 10)
		partitionName := p.partitionName(tableName, borders)

		stmt := "create table " + partitionName + " partition of " + tableName + " for values from (" + from + ") to (" + to + ")"
		log.Debug().Msg("Create partition '" + partitionName + "' with statement: " + stmt)

		_, err := p.driver.Db.Exec(stmt)
		if err != nil {
			panic(err)
		}
	}
}

func (p *PartitionManager) removeOld(tableName string, bordersList []partitionBorders) {
	allPartitions := p.allPartitions(tableName)
	actualPartitions := make([]string, len(bordersList))

	for i, borders := range bordersList {
		actualPartitions[i] = p.partitionName(tableName, borders)
	}

	oldPartitions := p.filterPartitionsGetOld(allPartitions, actualPartitions)

	for _, name := range oldPartitions {
		log.Debug().Msg("drop table " + name)

		_, err := p.driver.Db.Exec("drop table " + name)

		if err != nil {
			panic(err)
		}
	}
}

func (p *PartitionManager) allPartitions(tableName string) []string {
	rows, err := p.driver.Db.Query("SELECT child.relname FROM pg_inherits inherits, pg_class parent, pg_class child"+
		" WHERE inherits.inhparent = parent.oid and inherits.inhrelid = child.oid and parent.relname=$1", tableName)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	allPartitions := make([]string, 0)

	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}

		allPartitions = append(allPartitions, name)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return allPartitions
}

func (p *PartitionManager) filterPartitionsGetOld(all []string, actual []string) []string {
	old := make([]string, len(all)-len(actual))
	actualValues := set.SliceToSet(actual)
	i := 0

	for _, value := range all {
		if _, ok := actualValues[value]; !ok {
			old[i] = value
			i++
		}
	}

	return old
}

func (p *PartitionManager) actualBorders(tableName string, timeRange int, count int) []partitionBorders {
	now := time.Now().Unix()
	timeRangeInt64 := int64(timeRange)
	countInt64 := int64(count)
	partitionsSettings := make([]partitionBorders, count)

	for i := int64(0); i < countInt64; i++ {
		startValue, endValue := p.partitionBorders(now, timeRangeInt64, i)
		partitionsSettings[i] = partitionBorders{startValue: startValue, endValue: endValue}
	}

	return partitionsSettings
}

// get start and end time for the partition
// where start time is the current time truncated to a whole number of timeRange periods plus the offset.
func (p *PartitionManager) partitionBorders(now int64, timeRange int64, offset int64) (int64, int64) {
	start := now/timeRange*timeRange + timeRange*offset
	end := start + timeRange

	return start, end
}

func (p *PartitionManager) partitionName(tableName string, borders partitionBorders) string {
	if borders.startValue == borders.endValue {
		panic("Incorrect configuration for table '" + tableName + "' partitions. Start and end values are equal.")
	}

	from := strconv.FormatInt(borders.startValue, 10)
	to := strconv.FormatInt(borders.endValue, 10)

	return tableName + "_" + from + "_" + to
}

func (p *PartitionManager) nextEvictionTime() {
	// todo: calculate next eviction time
}

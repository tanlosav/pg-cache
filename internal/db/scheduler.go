package db

import (
	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type Scheduler struct {
	config           *configuration.Configuration
	partitionManager *PartitionManager
}

func NewScheduler(config *configuration.Configuration, partitionManager *PartitionManager) *Scheduler {
	return &Scheduler{
		config:           config,
		partitionManager: partitionManager,
	}
}

func (s *Scheduler) start() {
	// todo: calculate next eviction time and schedule the task
	for bucket, settings := range s.config.Cache.Buckets {
		for partitionNumber := 0; partitionNumber < settings.Sharding.PartitionsCount; partitionNumber++ {
			s.createBucket(bucket, settings.KeysCount, partitionNumber, settings.Eviction)
		}
	}

	nextEvictionTime := s.partitionManager.nextEvictionTime()
	go s.rotatePartitions()
}

func (s *Scheduler) rotatePartitions() {
	log.Info().Msg("Rotate partitions.")
}

package configuration

import (
	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/internal/cmd"
)

type ConfigurationSource interface {
	Configuration(source string) Configuration
}

type Configuration struct {
	Db     Db     `yaml:"db"`
	Cache  Cache  `yaml:"cache"`
	Server Server `yaml:"server"`
}

type Cache struct {
	Buckets map[string]Bucket `yaml:"buckets"`
	Hash    string            `yaml:"hash"`
}

type Db struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Bucket struct {
	KeysCount int `yaml:"keysCount"`
	Sharding  struct {
		PartitionsCount int `yaml:"partitionsCount"`
	} `yaml:"sharding"`
	Eviction Eviction `yaml:"eviction"`
}

type Eviction struct {
	Policy                   string `yaml:"policy"`
	PartitionTimeRange       int    `yaml:"partitionTimeRange"`
	ActualPartitionsCount    int    `yaml:"actualPartitionsCount"`
	RemainingPartitionsCount int    `yaml:"remainingPartitionsCount"`
}

type Server struct {
	Port int `yaml:"port"`
}

const (
	DEFAULT_PARTITIONS_COUNT           = 1
	DEFAULT_EVICTION_POLICY            = EVICTION_POLICY_NONE
	DEFAULT_PARTITION_TIME_RANGE       = 3600
	DEFAULT_ACTUAL_PARTITIONS_COUNT    = 2
	DEFAULT_REMAINING_PARTITIONS_COUNT = 1
	DEFAULT_SERVER_PORT                = 8080
)

const (
	EVICTION_POLICY_NONE     = "none"
	EVICTION_POLICY_DELETE   = "delete"
	EVICTION_POLICY_TRUNCATE = "truncate"
)

func NewConfiguration(opts *cmd.CmdLineOpts) *Configuration {
	var configProvider ConfigurationSource

	switch opts.ConfigurationProvider {
	case "file":
		configProvider = NewFileSource()
	default:
		panic("Unsupported configuration provider")
	}

	config := configProvider.Configuration(opts.ConfigurationSource)

	log.Printf("Loaded configuration: %+v", config)

	return &config
}

func applyBucketDefaultValues(conf *Configuration) {
	for bucket, opts := range conf.Cache.Buckets {
		if opts.Sharding.PartitionsCount == 0 {
			opts.Sharding.PartitionsCount = DEFAULT_PARTITIONS_COUNT
		}

		if opts.Eviction.Policy == "" {
			opts.Eviction.Policy = DEFAULT_EVICTION_POLICY
		}

		if opts.Eviction.PartitionTimeRange == 0 {
			opts.Eviction.PartitionTimeRange = DEFAULT_PARTITION_TIME_RANGE
		}

		if opts.Eviction.ActualPartitionsCount == 0 {
			opts.Eviction.ActualPartitionsCount = DEFAULT_ACTUAL_PARTITIONS_COUNT
		}

		if opts.Eviction.RemainingPartitionsCount == 0 {
			opts.Eviction.RemainingPartitionsCount = DEFAULT_REMAINING_PARTITIONS_COUNT
		}

		conf.Cache.Buckets[bucket] = opts
	}

	if conf.Server.Port == 0 {
		conf.Server.Port = DEFAULT_SERVER_PORT
	}
}

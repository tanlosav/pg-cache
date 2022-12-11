package configuration

import (
	"github.com/rs/zerolog/log"
	"github.com/tanlosav/pg-cache/internal/cmd"
)

type ConfigurationSource interface {
	Configuration(source string) Configuration
}

type Configuration struct {
	Db    Db    `yaml:"db"`
	Cache Cache `yaml:"cache"`
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
		Partition       string `yaml:"partitionName"`
		PartitionsCount int    `yaml:"partitionsCount"`
	} `yaml:"sharding"`
}

const (
	DEFAULT_PARTITIONS_COUNT = 100
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

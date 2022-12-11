package configuration

import (
	"os"

	"github.com/go-yaml/yaml"
	"github.com/rs/zerolog/log"
)

type FileSource struct {
	initialized bool
	config      Configuration
}

func NewFileSource() *FileSource {
	return &FileSource{}
}

func (conf *FileSource) Configuration(source string) Configuration {
	if !conf.initialized {
		conf.load(source)
	}

	return conf.config
}

func (conf *FileSource) load(source string) {
	log.Info().Msg("Load configuration from: " + source)

	file, err := os.Open(source)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&conf.config); err != nil {
		panic(err)
	}

	conf.applyBucketDefaultValues()

	log.Info().Msg("Configuration loaded")

	conf.initialized = true
}

func (conf *FileSource) applyBucketDefaultValues() {
	for bucket, opts := range conf.config.Cache.Buckets {
		if opts.Sharding.PartitionsCount == 0 {
			opts.Sharding.PartitionsCount = DEFAULT_PARTITIONS_COUNT
		}

		conf.config.Cache.Buckets[bucket] = opts
	}
}

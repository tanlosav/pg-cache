package configuration

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
		Partition           string `yaml:"partitionName"`
		PartitionsCount     int    `yaml:"partitionsCount"`
		VirtualBucketsCount int    `yaml:"virtualBucketsCount"`
	} `yaml:"sharding"`
}

const (
	DEFAULT_PARTITIONS_COUNT      = 100
	DEFAULT_VIRTUAL_BUCKETS_COUNT = 1024
)

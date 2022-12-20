package partition

import (
	"crypto/md5"
	"fmt"
	"hash"
	"math/big"

	"github.com/spaolacci/murmur3"
	"github.com/tanlosav/pg-cache/internal/configuration"
)

type PartitionManager struct {
	config configuration.Configuration
	digest hash.Hash
}

func NewPartitionManager(config configuration.Configuration) *PartitionManager {
	return &PartitionManager{
		config: config,
		digest: getDigest(config),
	}
}

// todo: calculate instance id
func (pm *PartitionManager) GetPartition(bucket string, key string) string {
	bucketSettings := pm.config.Cache.Buckets[bucket]
	partitionsCount := bucketSettings.Sharding.PartitionsCount
	partitionNumber := pm.getPartitionNumber(key, partitionsCount)

	return fmt.Sprintf("%s_%d", bucket, partitionNumber)
}

func (pm *PartitionManager) getPartitionNumber(key string, partitionsCount int) int64 {
	hash := pm.getHash(key)

	return pm.calculatePartitionNumber(hash, partitionsCount)
}

func (pm *PartitionManager) getHash(value string) []byte {
	pm.digest.Reset()
	pm.digest.Write([]byte(value))

	return pm.digest.Sum(nil)
}

func (pm *PartitionManager) calculatePartitionNumber(hash []byte, partitionsCount int) int64 {
	a := new(big.Int)
	a.SetString(string(hash), 10)

	b := new(big.Int)
	b.SetInt64(int64(partitionsCount))

	modulo := new(big.Int)

	return modulo.Mod(a, b).Int64()
}

func getDigest(config configuration.Configuration) hash.Hash {
	switch config.Cache.Hash {
	case "md5":
		return md5.New()
	case "murmur3":
		return murmur3.New32()
	default:
		panic("Unsupported hash function: " + config.Cache.Hash)
	}
}

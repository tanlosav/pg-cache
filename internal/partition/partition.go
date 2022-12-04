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

func (pm *PartitionManager) GetPartition(bucket string, key string) string {
	bucketSettings := pm.config.Cache.Buckets[bucket]
	partitionName := bucketSettings.Sharding.Partition
	vbCount := bucketSettings.Sharding.VirtualBucketsCount
	partitionsCount := bucketSettings.Sharding.PartitionsCount
	partitionNumber := pm.getPartitionNumber(key, vbCount, partitionsCount)

	return fmt.Sprintf("%s_%d", partitionName, partitionNumber)
}

func (pm *PartitionManager) getPartitionNumber(key string, vbCount int, partitionsCount int) int64 {
	hash := pm.getHash(key)
	vbNumber := pm.calculateVirtualBucketNumber(hash, vbCount)

	return pm.calculatePartitionNumber(vbNumber, partitionsCount)
}

func (pm *PartitionManager) getHash(value string) []byte {
	pm.digest.Reset()
	pm.digest.Write([]byte(value))

	return pm.digest.Sum(nil)
}

func (pm *PartitionManager) calculateVirtualBucketNumber(hash []byte, vbCount int) int64 {
	a := new(big.Int)
	a.SetString(string(hash), 10)

	b := new(big.Int)
	b.SetInt64(int64(vbCount))

	modulo := new(big.Int)

	return modulo.Mod(a, b).Int64()
}

func (pm *PartitionManager) calculatePartitionNumber(vbNumber int64, partitionsCount int) int64 {
	a := new(big.Int)
	a.SetInt64(vbNumber)

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

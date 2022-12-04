package cache

type Cache interface {
	Get(bucket string, key string) string
	Create(bucket string, key string, document string, ttl int)
	Update(bucket string, key string, document string, ttl int)
	Delete(bucket string, key string)
	Touch(bucket string, key string, ttl int)
	Clean(bucket string)
}

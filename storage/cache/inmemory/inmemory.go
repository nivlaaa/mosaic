package inmemory

import ()

const (
	cacheName = "inmemory"
)

type cache struct {
	buckets map[string][]string
}

func New() *cache {
	c := &cache{}
	c.buckets = make(map[string][]string)
	return c
}

func (c *cache) Name() string {
	return cacheName
}

func (c *cache) Get(bucket string) ([]string, bool) {
	items, cached := c.buckets[bucket]
	return items, cached
}

func (c *cache) Add(item, bucket string) {
	c.buckets[bucket] = append(c.buckets[bucket], item)
}

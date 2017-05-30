package cache

import ()

type Cache interface {
	// Returns the name of the driver
	Name() string

	// Returns a list of all the items in a bucket
	Get(bucket string) ([]string, bool)

	// Add a item to the bucket
	Add(item, bucket string)
}

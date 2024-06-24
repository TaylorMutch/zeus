package auth

// TODO - setup an in-memory authentication cache which syncs from a storage provider

// TODO - use concurrent-map to map keys to tenants
type Cache interface {
	Get(key string) (Credential, error)
	Refresh() error
}

type cacheImpl struct {
	// Add fields here, such as memcached client details
}

func NewCache() Cache {
	return &cacheImpl{}
}

func (c *cacheImpl) Get(key string) (Credential, error) {
	// Implement this method
	return Credential{}, nil
}

func (c *cacheImpl) Refresh() error {
	// Implement this method
	return nil
}

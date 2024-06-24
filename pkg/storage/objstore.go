package storage

import (
	"os"

	gokitlog "github.com/go-kit/log"
	"github.com/thanos-io/objstore"
	objstoreclient "github.com/thanos-io/objstore/client"
)

// ObjectStore is an interface for storing and retrieving objects
type ObjectStore interface {
	// GetObject retrieves an object from the object store
	GetObject(key string) ([]byte, error)
	// ListObjects lists all objects in the object store with the given prefix
	ListObjects(prefix string) ([]string, error)
	// PutObject stores an object in the object store
	PutObject(key string, data []byte) error
	// DeleteObject removes an object from the object store
	DeleteObject(key string) error
}

type objectStoreImpl struct {
	// Add fields here
	b objstore.Bucket
}

func NewObjectStore(componentName string, bucketConf []byte) (ObjectStore, error) {
	s := objectStoreImpl{}
	logger := gokitlog.NewJSONLogger(gokitlog.NewSyncWriter(os.Stderr))
	b, err := objstoreclient.NewBucket(
		logger,
		bucketConf,
		componentName,
	)
	if err != nil {
		//logger.Log("msg", "failed to create object store bucket", "err", err)
		return nil, err
	}
	s.b = b

	return &s, nil
}

func (s *objectStoreImpl) GetObject(key string) ([]byte, error) {
	// Implement this method
	return nil, nil
}

func (s *objectStoreImpl) ListObjects(prefix string) ([]string, error) {
	// Implement this method
	return nil, nil
}

func (s *objectStoreImpl) PutObject(key string, data []byte) error {
	// Implement this method
	return nil
}

func (s *objectStoreImpl) DeleteObject(key string) error {
	// Implement this method
	return nil
}

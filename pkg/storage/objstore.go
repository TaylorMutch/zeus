package storage

import (
	"bytes"
	"context"
	"os"

	gokitlog "github.com/go-kit/log"
	"github.com/thanos-io/objstore"
	objstoreclient "github.com/thanos-io/objstore/client"
)

// ObjectStore is an interface for storing and retrieving objects
type ObjectStore interface {
	// Exists checks if an object exists in the object store
	Exists(ctx context.Context, key string) (bool, error)
	// GetObject retrieves an object from the object store
	GetObject(ctx context.Context, key string) ([]byte, error)
	// PutObject stores an object in the object store
	PutObject(ctx context.Context, key string, data []byte) error
	// DeleteObject removes an object from the object store
	DeleteObject(ctx context.Context, key string) error
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
		logger.Log("msg", "failed to create object store bucket", "err", err)
		return nil, err
	}
	s.b = b

	return &s, nil
}

func (s *objectStoreImpl) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := s.b.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *objectStoreImpl) GetObject(ctx context.Context, key string) ([]byte, error) {
	reader, err := s.b.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	data := buf.Bytes()

	return data, nil
}

func (s *objectStoreImpl) PutObject(ctx context.Context, key string, data []byte) error {
	return s.b.Upload(ctx, key, bytes.NewReader(data))
}

func (s *objectStoreImpl) DeleteObject(ctx context.Context, key string) error {
	return s.b.Delete(ctx, key)
}

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/TaylorMutch/zeus/pkg/storage"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const (
	CredentialStoreStoragePrefix = "credentials"
)

var (
	CredentialDoesNotExistError = errors.New("credential does not exist")
)

type CredentialStore interface {
	// GetCredential retrieves a credential from the store
	GetCredential(ctx context.Context, id string) (*Credential, error)

	// CacheCredential caches a credential in the store
	CacheCredential(id string, cred *Credential)
}

type ObjectCredentialStore struct {
	store storage.ObjectStore
	cache cmap.ConcurrentMap[string, *Credential]
}

func NewObjectCredentialStore(store storage.ObjectStore) (CredentialStore, error) {
	// TODO - start periodically refreshing the store
	return &ObjectCredentialStore{
		store: store,
		cache: cmap.New[*Credential](),
	}, nil
}

func (ocs *ObjectCredentialStore) GetCredential(ctx context.Context, id string) (*Credential, error) {
	// Check if the credential is in the cache
	cred, ok := ocs.cache.Get(id)
	if ok {
		return cred, nil
	}

	// If not, check the store
	objectLocation := fmt.Sprintf("%s/%s", CredentialStoreStoragePrefix, id)
	exists, err := ocs.store.Exists(ctx, objectLocation)
	if err != nil {
		slog.Error("failed to check if credential exists", "id", id, "error", err)
		return nil, err
	}
	if !exists {
		slog.Info("credential does not exist", "id", id, "error", err)
		return nil, CredentialDoesNotExistError
	}
	byts, err := ocs.store.GetObject(ctx, objectLocation)
	if err != nil {
		slog.Error("failed to get credential", "id", id, "error", err)
		return nil, err
	}

	// Unmarshal the credential
	result := new(Credential)
	err = json.Unmarshal(byts, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ocs *ObjectCredentialStore) CacheCredential(id string, cred *Credential) {
	ocs.cache.Set(id, cred)
}

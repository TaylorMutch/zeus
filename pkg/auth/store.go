package auth

import (
	"encoding/json"
	"fmt"

	"github.com/TaylorMutch/zeus/pkg/storage"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const (
	CredentialStoreStoragePrefix = "credentials"
)

type CredentialStore interface {
	// GetCredential retrieves a credential from the store
	GetCredential(id string) (*Credential, error)
}

type ObjectCredentialStore struct {
	store storage.ObjectStore
	cache cmap.ConcurrentMap[string, Credential]
}

func NewObjectCredentialStore(store storage.ObjectStore) (CredentialStore, error) {
	// TODO - start periodically refreshing the store
	return &ObjectCredentialStore{
		store: store,
		cache: cmap.New[Credential](),
	}, nil
}

func (ocs *ObjectCredentialStore) GetCredential(id string) (*Credential, error) {
	// Implement this method
	// Check if the credential is in the cache

	cred, ok := ocs.cache.Get(id)
	if ok {
		return &cred, nil
	}

	// If not, check the store
	objectLocation := fmt.Sprintf("%s/%s", CredentialStoreStoragePrefix, id)
	byts, err := ocs.store.GetObject(objectLocation)
	if err != nil {
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

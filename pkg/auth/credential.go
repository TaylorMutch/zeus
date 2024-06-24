package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
)

const (
	// All credential password prefixes
	CredentialPasswordPrefix = "zeus_"

	// Raw passwords should be 24 characters long
	PasswordLength = 24

	// Allowed characters for random string generation
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

/*
	From a credential we should know the following:
	- The tenant ID
*/

// TODO - should this be a proto?
// Credential represents a tenant's credential.
// It should contain the tenant's ID, username, and password.
// It can be used as a basic authentication token which can be used for most requests.
type Credential struct {
	// ID is an internal ID of the token
	ID string

	// TenantID is the tenant which the credential belongs to.
	TenantID string

	// Username is the username used for authorizing traffic for the tenant's credential.
	// It should be a human readable string.
	Username string

	// Password is the HashedPassword used for authorizing traffic for
	// the tenant's credential.
	// Raw passwords should be 24 characters long and contain alphanumeric characters.
	// Raw passwords are never stored, only the hash and salt are stored.
	Password HashedPassword
}

// CreateToken takes a credential and returns a token
// than can be used for authentication.
// It should generate a base64 encoded token from the credential
// prefixed with `zeus_`.
func (c *Credential) CreateToken() string {
	return ""

}

func newCredentialID() string {
	// TODO - implement this
	s, _ := GenerateRandomString(24)
	return s
}

func hashCredential(t Credential) string {
	// TODO - implement this
	return ""
}

func init() {
	assertAvailablePRNG()
}

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {

	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

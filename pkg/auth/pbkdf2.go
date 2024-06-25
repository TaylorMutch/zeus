package auth

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Raw passwords should be 24 characters long
	DefaultPasswordLength = 24

	// MinSaltSize a minimum salt size recommended by the RFC
	SaltSize = 64

	// Iterations
	Iterations = 10000
)

// DefaultPasswordFactory is the default password factory for zeus
var DefaultPasswordFactory = NewPasswordFactory()

type PasswordFactory struct {
	Digest     func() hash.Hash
	SaltSize   int
	KeyLen     int
	Iterations int
}

type HashedPassword struct {
	CipherText string `json:"cipher_text"`
	Salt       string `json:"salt"`
}

func NewPasswordFactory() *PasswordFactory {
	return &PasswordFactory{
		Digest:     crypto.SHA512.New,
		SaltSize:   SaltSize,
		KeyLen:     DefaultPasswordLength,
		Iterations: Iterations,
	}
}

func (p *PasswordFactory) generateSalt() string {
	saltBytes := make([]byte, p.SaltSize)
	_, err := rand.Read(saltBytes)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(saltBytes)
}

func (p *PasswordFactory) HashPassword(password string) HashedPassword {
	saltString := p.generateSalt()
	salt := bytes.NewBufferString(saltString).Bytes()
	df := pbkdf2.Key([]byte(password), salt, p.Iterations, p.KeyLen, p.Digest)
	cipherText := base64.StdEncoding.EncodeToString(df)
	return HashedPassword{CipherText: cipherText, Salt: saltString}
}

func (p *PasswordFactory) VerifyPassword(password, cipherText, salt string) bool {
	saltBytes := bytes.NewBufferString(salt).Bytes()
	df := pbkdf2.Key([]byte(password), saltBytes, p.Iterations, p.KeyLen, p.Digest)
	return equal(cipherText, df)
}

// check per bit by applying bitwise XOR
// first, decode the base64 string to bytes
// for example
// 114  1110010
// 114  1110010
// ----------------- xor
//
//	0000000
func equal(cipherText string, newCipherText []byte) bool {
	x, _ := base64.StdEncoding.DecodeString(cipherText)
	diff := uint64(len(x)) ^ uint64(len(newCipherText))

	for i := 0; i < len(x) && i < len(newCipherText); i++ {
		diff |= uint64(x[i]) ^ uint64(newCipherText[i])
	}

	return diff == 0
}

package encryptutil

import (
	"crypto/cipher"
)

const (
	saltSize   = 8
	iterations = 32
)

type Encryptor interface {
	Hash(plaintext string) (string, error)
	CompareHashedWithPlainText(hashedPlaintext, plaintext string) (bool, error)

	Encrypt(plaintext string, isDeterministic bool) (string, error)
	EncryptWithMarshalling(v interface{}, isDeterministic bool) (string, error)

	Decrypt(ciphertext string) (string, error)
	DecryptWithUnmarshalling(ciphertext string, v interface{}) error
}

type ByteEncryptor interface {
	Encrypt(plaintext []byte, isDeterministic bool) ([]byte, error)
	EncryptWithMarshalling(v interface{}, isDeterministic bool) ([]byte, error)

	Decrypt(ciphertext []byte) ([]byte, error)
	DecryptWithUnmarshalling(ciphertext []byte, v interface{}) error
}

type encryptorImpl struct {
	key string
	*byteEncryptorImpl
}

type byteEncryptorImpl struct {
	aead       cipher.AEAD
	saltSize   int
	iterations int
}

func (e *byteEncryptorImpl) SetSaltSize(set int) {
	e.saltSize = set
}

func (e *byteEncryptorImpl) SetIterations(set int) {
	e.iterations = set
}

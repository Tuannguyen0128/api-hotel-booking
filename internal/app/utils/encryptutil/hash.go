package encryptutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"

	"golang.org/x/crypto/pbkdf2"
)

func (e *encryptorImpl) CompareHashedWithPlainText(hashedPlaintext, plaintext string) (bool, error) {
	if hashedPlaintext == "" {
		return false, errors.New("hashed password is empty")
	}
	hashedBytes, err := hex.DecodeString(hashedPlaintext)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hashedBytes, e.hash([]byte(plaintext), hashedBytes[:e.saltSize])), nil
}

func (e *encryptorImpl) Hash(plaintext string) (string, error) {
	salt, err := e.generateRandomBytes(e.saltSize)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(e.hash([]byte(plaintext), salt)), nil
}

func (e *encryptorImpl) hash(pw, salt []byte) []byte {
	ret := make([]byte, len(salt))
	copy(ret, salt)
	return append(ret, pbkdf2.Key(pw, salt, e.iterations, sha256.Size, sha256.New)...)
}

func (e *encryptorImpl) generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

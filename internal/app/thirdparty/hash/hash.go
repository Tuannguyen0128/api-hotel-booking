package hash

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/pbkdf2"
)

var ErrEmptyPass = errors.New("hashed password is empty")

const (
	saltSize   = 8
	iterations = 32
)

func ComparePassword(hashedPw, pw string) (bool, error) {
	if hashedPw == "" {
		return false, ErrEmptyPass
	}
	hashedBytes, err := hex.DecodeString(hashedPw)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hashedBytes, hashPassword([]byte(pw), hashedBytes[:saltSize])), nil
}

func NewHashPassword(pw string) (string, error) {
	salt, err := generateRandomBytes(saltSize)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hashPassword([]byte(pw), salt)), nil
}

func hashPassword(pw, salt []byte) []byte {
	ret := make([]byte, len(salt))
	copy(ret, salt)
	return append(ret, pbkdf2.Key(pw, salt, iterations, sha256.Size, sha256.New)...)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

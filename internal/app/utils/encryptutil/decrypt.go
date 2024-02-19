package encryptutil

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

func (e *byteEncryptorImpl) DecryptWithUnmarshalling(ciphertext []byte, v interface{}) error {
	if plaintext, err := e.Decrypt(ciphertext); err != nil {
		return err

	} else if err := json.Unmarshal(plaintext, v); err != nil {
		return errors.New(fmt.Sprintf("unmarshal error:%s", err.Error()))
	} else {
		return nil
	}
}

func (e *byteEncryptorImpl) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := e.aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New(fmt.Sprintf(
			"decoded ciphertext has length %v which is less than nonce size %v",
			len(ciphertext), nonceSize))

	}

	if plaintext, err := e.aead.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil); err != nil {
		return nil,
			errors.New(fmt.Sprintf("cannot decrypt, invalid ciphertext error:%s", err.Error()))

	} else {
		return plaintext, nil
	}
}

func (e *encryptorImpl) DecryptWithUnmarshalling(ciphertext string, v interface{}) error {
	ciphertextDecoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot decrypt, cannot decode base64 to string error:%s", err.Error()))
	}
	return e.byteEncryptorImpl.DecryptWithUnmarshalling(ciphertextDecoded, v)
}

func (e *encryptorImpl) Decrypt(ciphertext string) (string, error) {
	ciphertextDecoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errors.New(fmt.Sprintf("cannot decrypt, cannot decode base64 to string error:%s", err.Error()))
	}

	plaintext, err := e.byteEncryptorImpl.Decrypt(ciphertextDecoded)
	if err != nil {
		return "", errors.New(fmt.Sprintf("cannot decrypt, invalid ciphertext error:%s", err.Error()))
	}

	return string(plaintext), nil
}

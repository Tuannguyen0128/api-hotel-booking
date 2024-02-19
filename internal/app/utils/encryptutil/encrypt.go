package encryptutil

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

func (e *byteEncryptorImpl) EncryptWithMarshalling(v interface{}, isDeterministic bool) ([]byte, error) {
	if bs, err := json.Marshal(v); err != nil {
		return nil, errors.New(fmt.Sprintf("marshal error:%s", err.Error()))
	} else {
		return e.Encrypt(bs, isDeterministic)
	}
}

func (e *byteEncryptorImpl) Encrypt(plaintext []byte, isDeterministic bool) ([]byte, error) {
	if nonce, err := e.getNonce(isDeterministic); err != nil {
		return nil, err

	} else {
		return e.aead.Seal(nonce, nonce, plaintext, nil), nil
	}
}

func (e *byteEncryptorImpl) getNonce(isDeterministic bool) ([]byte, error) {
	nonce := make([]byte, e.aead.NonceSize())

	if isDeterministic {
		for i := 0; i < len(nonce); i++ {
			nonce[i] = '0'
		}
		return nonce, nil

	} else if _, err := rand.Read(nonce); err != nil {
		return []byte{}, errors.New(fmt.Sprintf("encrypt error, cannot create nonce:%s", err.Error()))
	}

	return nonce, nil
}

func (e *encryptorImpl) EncryptWithMarshalling(v interface{}, isDeterministic bool) (string, error) {
	return e.convertToB64(e.byteEncryptorImpl.EncryptWithMarshalling(v, isDeterministic))
}

func (e *encryptorImpl) Encrypt(plaintext string, isDeterministic bool) (string, error) {
	return e.convertToB64(e.byteEncryptorImpl.Encrypt([]byte(plaintext), isDeterministic))
}

func (e *encryptorImpl) convertToB64(bs []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

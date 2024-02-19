package encryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"

	"api-hotel-booking/internal/app/thirdparty/logger"
)

var log = logger.WithModule("encryptor")

type util struct {
}

func NewUtil() *util {
	return &util{}
}

func (u *util) NewEncryptor(key string) Encryptor {
	bKey := GetKey32Bytes(key)
	return &encryptorImpl{
		key:               string(bKey),
		byteEncryptorImpl: interface{}(u.NewByteEncryptor(bKey)).(*byteEncryptorImpl),
	}
}

func (u *util) NewByteEncryptor(key []byte) ByteEncryptor {
	if len(key) != 32 {
		key = GetKey32Bytes(string(key))
	}

	if block, err := aes.NewCipher(key); err != nil {
		log.Panic("cannot create cipher with key:%s error:%s", string(key), err.Error())
		return nil

	} else if aead, err := cipher.NewGCM(block); err != nil {
		log.Panic("cannot create GCM with key:%s error:%s", string(key), err.Error())
		return nil

	} else {
		return &byteEncryptorImpl{
			aead: aead,
		}
	}
}

func GetKey32Bytes(keyString string) []byte {
	keyBytes := []byte(keyString)

	if len(keyBytes) != 32 {
		key32bytes := sha256.Sum256(keyBytes)

		return key32bytes[:]

	} else {
		return keyBytes
	}
}

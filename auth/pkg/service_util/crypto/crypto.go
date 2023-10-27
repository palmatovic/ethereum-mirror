package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
)

type Key struct {
	encryptionKey string
}

func NewKey(encryptionKey string) *Key {
	return &Key{encryptionKey: encryptionKey}
}

func (k *Key) Encrypt(decryptedText string) (*string, error) {
	aes256Cipher, err := aes.NewCipher([]byte(k.encryptionKey))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aes256Cipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	encryptedText := gcm.Seal(nonce, nonce, []byte(decryptedText), nil)
	encoded := hex.EncodeToString(encryptedText)
	return &encoded, nil
}

func (k *Key) Decrypt(encryptedText string) (*string, error) {
	aes256Cipher, err := aes.NewCipher([]byte(k.encryptionKey))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aes256Cipher)
	if err != nil {
		return nil, err
	}
	baEncryptedText, err := hex.DecodeString(encryptedText)
	if err != nil {
		return nil, err
	}
	if len(baEncryptedText) <= gcm.NonceSize() {
		return nil, errors.New("cannot decrypt text with invalid length")
	}
	plainText, err := gcm.Open(nil, baEncryptedText[:gcm.NonceSize()], baEncryptedText[gcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
	plainTexString := string(plainText)
	return &plainTexString, nil
}

func (k *Key) DecryptFilepath(filepath string) (*string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	decryptedText, err := k.Decrypt(string(file))
	if err != nil {
		return nil, err
	}
	return decryptedText, nil
}

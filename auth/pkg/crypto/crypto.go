package crypto

type Key struct {
	encryptionKey string
}

func NewKey(encryptionKey string) *Key {
	return &Key{encryptionKey: encryptionKey}
}

func (k *Key) Encrypt(decryptedText string) ([]byte, error) {

}

func (k *Key) Decrypt(encryptedText string) ([]byte, error) {

}

func (k *Key) DecryptFilepath(filepath string) ([]byte, error) {

}

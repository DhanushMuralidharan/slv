package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/json"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/curve25519"
	"gopkg.in/yaml.v3"
)

type encrypter struct {
	ephpublicKey *[]byte
	aead         *cipher.AEAD
}

type PublicKey struct {
	*keyBase
	encrypter *encrypter
}

func (publicKey *PublicKey) ShortId() *[]byte {
	sum := sha1.Sum(*publicKey.key)
	keyId := sum[len(sum)-shortKeyIdLength:]
	return &keyId
}

func (publicKey *PublicKey) Id() string {
	return commons.Encode(publicKey.toBytes())
}

func (publicKey PublicKey) MarshalYAML() (interface{}, error) {
	return publicKey.String(), nil
}

func (publicKey *PublicKey) UnmarshalYAML(value *yaml.Node) error {
	var pubKeyStr string
	if err := value.Decode(&pubKeyStr); err == nil {
		keyBase, err := kyeBaseFromString(pubKeyStr)
		if err != nil {
			return err
		}
		publicKey.keyBase = keyBase
	}
	return nil
}

func (publicKey PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(publicKey.String())
}

func (publicKey *PublicKey) UnmarshalJSON(data []byte) (err error) {
	var pubKeyStr string
	if err = json.Unmarshal(data, &pubKeyStr); err == nil {
		keyBase, err := kyeBaseFromString(pubKeyStr)
		if err != nil {
			return err
		}
		publicKey.keyBase = keyBase
	}
	return
}

func PublicKeyFromString(publicKeyStr string) (*PublicKey, error) {
	keyBase, err := kyeBaseFromString(publicKeyStr)
	if err != nil {
		return nil, err
	}
	return &PublicKey{
		keyBase: keyBase,
	}, nil
}

type SecretKey struct {
	*keyBase
	publicKey *PublicKey
}

func (secretKey *SecretKey) ShortId() *[]byte {
	pubKey, err := secretKey.PublicKey()
	if err != nil {
		return nil
	}
	return pubKey.ShortId()
}

func NewSecretKey(keyType KeyType) (secretKey *SecretKey, err error) {
	privKey := make([]byte, curve25519.ScalarSize)
	if _, err = rand.Read(privKey); err != nil {
		return nil, ErrGeneratingKey
	}
	return newSecretKey(&privKey, keyType), nil
}

func NewSecretKeyForPassword(password []byte, keyType KeyType) (secretKey *SecretKey, err error) {
	salt := make([]byte, argon2SaltLength)
	if _, err = rand.Read(salt); err != nil {
		return nil, ErrGeneratingKey
	}
	privKey := argon2.IDKey(password, salt, argon2Iterations,
		argon2Memory, argon2Threads, curve25519.ScalarSize)
	return newSecretKey(&privKey, keyType), nil
}

func newSecretKey(privKey *[]byte, keyType KeyType) *SecretKey {
	version := commons.Version
	return &SecretKey{
		keyBase: &keyBase{
			version: &version,
			public:  false,
			keyType: &keyType,
			key:     privKey,
		},
	}
}

func (secretKey *SecretKey) PublicKey() (*PublicKey, error) {
	if secretKey.publicKey == nil {
		key, err := curve25519.X25519(*secretKey.key, curve25519.Basepoint)
		if err != nil {
			return nil, ErrGeneratingKey
		}
		secretKey.publicKey = &PublicKey{
			keyBase: &keyBase{
				version: secretKey.version,
				public:  true,
				keyType: secretKey.keyType,
				key:     &key,
			},
		}
	}
	return secretKey.publicKey, nil
}

func SecretKeyFromString(secretKeyStr string) (*SecretKey, error) {
	keyBase, err := kyeBaseFromString(secretKeyStr)
	if err != nil {
		return nil, err
	}
	return &SecretKey{
		keyBase: keyBase,
	}, nil
}

package keymanagerv1

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

type KeyManager struct {
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	publicKeyPEM string
}

func New(privateKeyPEM []byte) (*KeyManager, error) {

	privateKey, err := generateRsaPrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	publicKeyPEM, err := generatePublicKeyPEM(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &KeyManager{
		privateKey:   privateKey,
		publicKey:    &privateKey.PublicKey,
		publicKeyPEM: publicKeyPEM,
	}, nil
}

func generateRsaPrivateKey(keyPEM []byte) (*rsa.PrivateKey, error) {
	keyPEM = bytes.TrimSpace(keyPEM)

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, errors.New("failed to parse pem block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {

		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}

		var ok bool
		if privateKey, ok = key.(*rsa.PrivateKey); !ok {
			return nil, errors.New("key is not rsa private key")
		}
	}

	return privateKey, nil
}

func generatePublicKeyPEM(publicKey *rsa.PublicKey) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

func (m *KeyManager) GetPrivateKey() *rsa.PrivateKey {
	return m.privateKey
}

func (m *KeyManager) GetPublicKey() *rsa.PublicKey {
	return m.publicKey
}

func (m *KeyManager) GetPublicKeyPEM() string {
	return m.publicKeyPEM
}

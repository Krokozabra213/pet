package keymanagerv1

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

type PublicManager struct {
	publicKey    *rsa.PublicKey
	publicKeyPEM string
}

func NewPublicManager(publicKeyPEM string) (PublicManager, error) {
	publicKey, err := ParsePublicKeyPEM(publicKeyPEM)
	if err != nil {
		return PublicManager{}, err
	}

	return PublicManager{
		publicKey:    publicKey,
		publicKeyPEM: publicKeyPEM,
	}, nil

}

func ParsePublicKeyPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	keyPEM := bytes.TrimSpace([]byte(publicKeyPEM))

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, errors.New("failed to parse pem block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not rsa public key")
	}

	return rsaPublicKey, nil
}

func (m *PublicManager) GetPublicKey() *rsa.PublicKey {
	return m.publicKey
}

func (m *PublicManager) GetPublicKeyPEM() string {
	return m.publicKeyPEM
}

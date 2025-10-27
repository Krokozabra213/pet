package keymanager

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

type KeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New(privateKeyPEM []byte) *KeyManager {

	privateKeyPEM = bytes.TrimSpace(privateKeyPEM)

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Fatal("failed to parse PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {

		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			log.Fatalf("failed to parse private key: %v", err)
		}

		var ok bool
		if privateKey, ok = key.(*rsa.PrivateKey); !ok {
			log.Fatal("key is not RSA private key")
		}
	}

	return &KeyManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

func (km *KeyManager) GetPublicKeyPEM() (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(km.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %v", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

func (km *KeyManager) GetPrivateKey() *rsa.PrivateKey {
	return km.privateKey
}

func (km *KeyManager) GetPublicKey() *rsa.PublicKey {
	return km.publicKey
}

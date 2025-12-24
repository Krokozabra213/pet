package keymanager

// import (
// 	"bytes"
// 	"crypto/rsa"
// 	"crypto/x509"
// 	"encoding/pem"
// 	"errors"
// )

// var (
// 	ErrDecodePem    = errors.New("failed to decode PEM block containing public key")
// 	ErrParsePem     = errors.New("failed to parse PEM block containing public key")
// 	ErrRSAPublicKey = errors.New("key is not RSA public key")
// )

// type PublicKeyManager struct {
// 	PublicKey *rsa.PublicKey
// }

// func NewPublic(publicKeyPEM []byte) (*PublicKeyManager, error) {
// 	publicKeyPEM = bytes.TrimSpace(publicKeyPEM)

// 	block, _ := pem.Decode(publicKeyPEM)
// 	if block == nil {
// 		return nil, ErrParsePem
// 	}

// 	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
// 	if err != nil {

// 		key, err := x509.ParsePKIXPublicKey(block.Bytes)
// 		if err != nil {

// 			return nil, ErrParsePem
// 		}

// 		var ok bool
// 		publicKey, ok = key.(*rsa.PublicKey)
// 		if !ok {
// 			return nil, ErrRSAPublicKey
// 		}
// 	}

// 	return &PublicKeyManager{
// 		PublicKey: publicKey,
// 	}, nil
// }

// func (manager *PublicKeyManager) GetPublicKey() *rsa.PublicKey {
// 	return manager.PublicKey
// }

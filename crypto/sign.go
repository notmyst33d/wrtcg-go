package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
)

func PrivateKeyFromBase64(key string) (ed25519.PrivateKey, error) {
	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return ed25519.PrivateKey(nil), err
	}
	return ed25519.PrivateKey(data), nil
}

func PublicKeyFromBase64(key string) (ed25519.PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return ed25519.PublicKey(nil), err
	}
	return ed25519.PublicKey(data), nil
}

func SignAnswer() {
}

package option

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type ConfigKeys struct {
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

type _Config struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	Signaling  SignalingOptions `json:"signaling"`
}

type Config _Config

func (o *Config) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_Config)(o))
	if err != nil {
		return err
	}
	keys := ConfigKeys{}
	err = json.Unmarshal(bytes, &keys)
	if err != nil {
		return err
	}
	if keys.PublicKey != "" {
		publicKey, err := base64.StdEncoding.DecodeString(keys.PublicKey)
		if err != nil {
			return err
		}
		o.PublicKey = ed25519.PublicKey(publicKey)
	}
	if keys.PrivateKey != "" {
		privateKey, err := base64.StdEncoding.DecodeString(keys.PrivateKey)
		if err != nil {
			return err
		}
		o.PrivateKey = ed25519.PrivateKey(privateKey)
	}
	return nil
}

type _SignalingOptions struct {
	Type            string                   `json:"type"`
	TelemostOptions TelemostSignalingOptions `json:"-"`
}

type SignalingOptions _SignalingOptions

func (o *SignalingOptions) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, (*_SignalingOptions)(o))
	if err != nil {
		return err
	}
	var v any
	switch o.Type {
	case "telemost":
		v = &o.TelemostOptions
	default:
		return errors.New("unknown signaling type: " + o.Type)
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return nil
}

type TelemostSignalingOptions struct {
	Link string `json:"link"`
}

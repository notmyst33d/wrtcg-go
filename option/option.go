package option

import (
	"encoding/json"
	"errors"
)

type Config struct {
	Signaling SignalingOptions `json:"signaling"`
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

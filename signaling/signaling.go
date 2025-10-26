package signaling

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var TokenEncoding = base64.URLEncoding

const TokenTypeOffer = "offer"
const TokenTypeAnswer = "answer"

type SignalingClient interface {
	GetTokenChannel() chan Token
	SendToken(token Token) error
	GetIceServers() ([]ICEServer, error)
}

type Token struct {
	RequestID uint64
	SDP       string
	TokenType string
	Signature string
}

func (t *Token) Decode(token string) error {
	segments := strings.Split(token, ".")
	if len(segments) == 1 {
		return errors.New("not a token")
	}
	signed := strings.HasSuffix(segments[0], "-signed")

	t.TokenType = strings.TrimSuffix(strings.TrimPrefix(segments[0], "wrtcg-"), "-signed")

	requestId, err := strconv.ParseInt(segments[1], 10, 64)
	if err != nil {
		return err
	}
	t.RequestID = uint64(requestId)

	sdp, err := TokenEncoding.DecodeString(segments[2])
	if err != nil {
		return err
	}
	t.SDP = string(sdp)

	if signed {
		signature, err := TokenEncoding.DecodeString(segments[3])
		if err != nil {
			return err
		}
		t.Signature = string(signature)
	}

	return nil
}

func (t *Token) Encode() string {
	return fmt.Sprintf("wrtcg-%s.%d.%s", t.TokenType, t.RequestID, TokenEncoding.EncodeToString([]byte(t.SDP)))
}

func (t *Token) EncodeSigned(privateKey ed25519.PrivateKey) string {
	signature := TokenEncoding.EncodeToString(ed25519.Sign(privateKey, []byte(t.SDP)))
	return fmt.Sprintf("wrtcg-%s-signed.%d.%s.%s", t.TokenType, t.RequestID, TokenEncoding.EncodeToString([]byte(t.SDP)), signature)
}

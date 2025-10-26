package signaling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/notmyst33d/wrtcg-go/v2/option"
)

var _ SignalingClient = (*TelemostSignaling)(nil)

type TelemostSignaling struct {
	Options            option.TelemostSignalingOptions
	Conn               *websocket.Conn
	ConnMutex          sync.Mutex
	ServerHelloChannel chan struct{}
	ServerHello        *ServerHelloMessage
	TokenChannel       chan Token
}

type ClientConfiguration struct {
	MediaServerURL string `json:"media_server_url"`
}

type ConnectionResponse struct {
	PeerID              string              `json:"peer_id"`
	RoomID              string              `json:"room_id"`
	Credentials         string              `json:"credentials"`
	ClientConfiguration ClientConfiguration `json:"client_configuration"`
}

type ParticipantMeta struct {
	Name        string `json:"name"`
	Role        string `json:"role"`
	Description string `json:"description"`
	SendAudio   bool   `json:"sendAudio"`
	SendVideo   bool   `json:"sendVideo"`
}

type RTCConfiguration struct {
	ICEServers []ICEServer `json:"iceServers"`
}

type Description struct {
	Meta           ParticipantMeta `json:"meta"`
	DisconnectedAt *uint64         `json:"disconnectedAt,omitempty"`
}

type Status struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type WSMessage struct {
	UID               string                    `json:"uid"`
	Hello             *HelloMessage             `json:"hello,omitempty"`
	ServerHello       *ServerHelloMessage       `json:"serverHello,omitempty"`
	Ping              *EmptyMessage             `json:"ping,omitempty"`
	ACK               *ACKMessage               `json:"ack,omitempty"`
	UpdateMe          *UpdateMeMessage          `json:"updateMe,omitempty"`
	UpsertDescription *UpsertDescriptionMessage `json:"upsertDescription,omitempty"`
}

type HelloMessage struct {
	ParticipantMeta        ParticipantMeta `json:"participantMeta"`
	ParticipantID          string          `json:"participantId"`
	RoomID                 string          `json:"roomId"`
	ServiceName            string          `json:"serviceName"`
	Credentials            string          `json:"credentials"`
	SDKInitializationID    string          `json:"sdkInitializationId"`
	DisablePublisher       bool            `json:"disablePublisher"`
	DisableSubscriber      bool            `json:"disableSubscriber"`
	DisableSubscriberAudio bool            `json:"disableSubscriberAudio"`
}

type ServerHelloMessage struct {
	RTCConfiguration RTCConfiguration `json:"rtcConfiguration"`
}

type UpdateMeMessage struct {
	ParticipantMeta ParticipantMeta `json:"participantMeta"`
}

type UpsertDescriptionMessage struct {
	Description []Description `json:"description"`
}

type ACKMessage struct {
	Status Status `json:"status"`
}

type EmptyMessage struct {
}

func NewTelemostSignaling(options option.TelemostSignalingOptions) TelemostSignaling {
	return TelemostSignaling{
		Options:            options,
		ServerHelloChannel: make(chan struct{}),
		TokenChannel:       make(chan Token),
	}
}

func (c *TelemostSignaling) writeJson(v any) error {
	c.ConnMutex.Lock()
	defer c.ConnMutex.Unlock()
	return c.Conn.WriteJSON(v)
}

func (c *TelemostSignaling) onMessage() {
	for {
		message := WSMessage{}
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			return
		}

		serverHello := false
		if message.ServerHello != nil {
			serverHello = true
		} else if message.UpsertDescription != nil {
			for _, desc := range message.UpsertDescription.Description {
				if desc.DisconnectedAt != nil {
					continue
				}
				rawToken := desc.Meta.Name
				if rawToken == "wrtcg-server" || rawToken == "wrtcg-client" {
					continue
				}
				token := Token{}
				err := token.Decode(rawToken)
				if err == nil {
					c.TokenChannel <- token
				} else {
					fmt.Println(err)
				}
			}
		}

		if message.ACK == nil {
			err = c.writeJson(WSMessage{
				UID: message.UID,
				ACK: &ACKMessage{
					Status: Status{
						Code:        "OK",
						Description: "",
					},
				},
			})
			if err != nil {
				return
			}
			if serverHello {
				c.ServerHello = message.ServerHello
				c.ServerHelloChannel <- struct{}{}
				go c.ping()
			}
		}
	}
}

func (c *TelemostSignaling) ping() {
	for {
		err := c.writeJson(WSMessage{
			UID:  uuid.NewString(),
			Ping: &EmptyMessage{},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(3 * time.Second)
	}
}

func (c *TelemostSignaling) Dial(name string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://cloud-api.yandex.ru/telemost_front/v2/telemost/conferences/%s/connection", url.QueryEscape(c.Options.Link)), nil)
	if err != nil {
		return err
	}
	req.Header.Set("client-instance-id", uuid.NewString())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	data := ConnectionResponse{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(data.ClientConfiguration.MediaServerURL, nil)
	if err != nil {
		return err
	}
	c.Conn = conn

	go c.onMessage()

	err = c.writeJson(WSMessage{
		UID: uuid.NewString(),
		Hello: &HelloMessage{
			ParticipantMeta: ParticipantMeta{
				Name:        name,
				Role:        "SPEAKER",
				Description: "",
				SendAudio:   false,
				SendVideo:   false,
			},
			ParticipantID:          data.PeerID,
			RoomID:                 data.RoomID,
			ServiceName:            "telemost",
			Credentials:            data.Credentials,
			SDKInitializationID:    uuid.NewString(),
			DisablePublisher:       true,
			DisableSubscriber:      true,
			DisableSubscriberAudio: true,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *TelemostSignaling) GetTokenChannel() chan Token {
	return c.TokenChannel
}

func (c *TelemostSignaling) SendToken(token Token) error {
	return c.writeJson(WSMessage{
		UID: uuid.NewString(),
		UpdateMe: &UpdateMeMessage{
			ParticipantMeta: ParticipantMeta{
				Name:        token.Encode(),
				Role:        "SPEAKER",
				Description: "",
				SendAudio:   false,
				SendVideo:   false,
			},
		},
	})
}

func (c *TelemostSignaling) GetIceServers() ([]ICEServer, error) {
	return c.ServerHello.RTCConfiguration.ICEServers, nil
}

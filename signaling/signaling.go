package signaling

type SignalingClient interface {
	ReceiveOffers() ([]string, error)
	SendAnswer(oid uint64, answer string) error
	GetIceServers() ([]ICEServer, error)
}

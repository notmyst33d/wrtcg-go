package signaling

type ICEServer struct {
	Urls       []string `json:"urls"`
	Username   string   `json:"username"`
	Credential string   `json:"credential"`
	Signature  *string  `json:"signature,omitempty"`
}

type SignalingPayload struct {
	Offer  string `json:"offer"`
	Answer string `json:"answer"`
}

type PollingData struct {
	IceServers []ICEServer        `json:"ice_servers"`
	Signaling  []SignalingPayload `json:"signaling"`
}

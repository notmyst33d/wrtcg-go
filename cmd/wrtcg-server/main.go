package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/notmyst33d/wrtcg-go/v2/option"
	"github.com/notmyst33d/wrtcg-go/v2/signaling"
)

func main() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("cannot open config:", err)
		os.Exit(1)
	}

	config := option.Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("cannot parse config:", err)
		os.Exit(1)
	}

	if config.Signaling.Type == "telemost" {
		client := signaling.NewTelemostSignaling(config.Signaling.TelemostOptions)
		err = client.Dial("wrtcg")
		if err != nil {
			fmt.Println("cannot start signaling client:", err)
		}
		<-client.ServerHelloChan
		for {
			fmt.Println(client.ReceiveOffers())
			time.Sleep(time.Second)
		}
	}
}

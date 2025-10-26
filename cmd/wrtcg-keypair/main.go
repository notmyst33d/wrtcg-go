package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Public Key: %s\n", base64.StdEncoding.EncodeToString(publicKey))
	fmt.Printf("Private Key: %s\n", base64.StdEncoding.EncodeToString(privateKey[:32]))
}

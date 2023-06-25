package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"encoding/json"
	"bytes"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	var oauthConfig = &clientcredentials.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		TokenURL:     "https://api.tailscale.com/api/v2/oauth/token",
	}

	client := oauthConfig.Client(context.Background())
	// Replace example.com with your tailnet name.
	// resp, err := client.Get("https://api.tailscale.com/api/v2/tailnet/vungle.com/devices")
	// if err != nil {
	// 	log.Fatalf("error getting keys: %v", err)
	// }

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("error reading response body: %v", err)
	// }
	
	// fmt.Printf("response: %s", string(body))
	
	// Call api to get all the keys
	getKeys_Resp, getKeysErr := client.Get("https://api.tailscale.com/api/v2/tailnet/vungle.com/keys")
	if getKeysErr != nil {
		log.Fatalf("error getting keys: %v", getKeysErr)
	}

	getKeysBody, getKeysErr := ioutil.ReadAll(getKeys_Resp.Body)
	if getKeysErr != nil {
		log.Fatalf("error reading response body: %v", getKeysErr)
	}
	
	fmt.Printf("get_keys_body response: %s", string(getKeysBody))	
	
	// Call api to create the key /api/v2/tailnet/{tailnet}/keys
	data := map[string]interface{}{
		"capabilities": map[string]interface{} {
			"devices": map[string]interface{} {
				"create": map[string]interface{} {
					"reusable": false,
					"ephemeral": true,
					"preauthorized": false,
					"tags": []string{"tag:prod"},
				},
			},
		},
		"expirySeconds": 3600,
		"description": "short description of key purpose",
	}	

	jsonData, formatErr := json.Marshal(data)
	if formatErr != nil {
		log.Fatalf("error creating body: %v", formatErr)
	}
	
	reqBody := bytes.NewBuffer(jsonData)

	createKeysResp, createKeysErr := client.Post("https://api.tailscale.com/api/v2/tailnet/vungle.com/keys", "application/json", reqBody)
	if getKeysErr != nil {
		log.Fatalf("error creating keys: %v", createKeysErr)
	}

	createKeysBody, createKeysErr := ioutil.ReadAll(createKeysResp.Body)
	if getKeysErr != nil {
		log.Fatalf("error reading response body: %v", createKeysErr)
	}
	
	fmt.Printf("create_keys_body response: %s", string(createKeysBody))		
	
}
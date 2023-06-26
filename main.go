package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"golang.org/x/oauth2/clientcredentials"
	// "encoding/json"
	// "bytes"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
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
	getKeysResp, getKeysErr := client.Get("https://api.tailscale.com/api/v2/tailnet/vungle.com/keys")
	if getKeysErr != nil {
		log.Fatalf("error getting keys: %v", getKeysErr)
	}

	getKeysBody, getKeysErr := ioutil.ReadAll(getKeysResp.Body)
	if getKeysErr != nil {
		log.Fatalf("error reading response body: %v", getKeysErr)
	}
	
	fmt.Printf("get_keys_body response: %s", string(getKeysBody))
	
	// Call api to create the key /api/v2/tailnet/{tailnet}/keys
	// data := map[string]interface{}{
	// 	"capabilities": map[string]interface{} {
	// 		"devices": map[string]interface{} {
	// 			"create": map[string]interface{} {
	// 				"reusable": true,
	// 				"ephemeral": true,
	// 				"preauthorized": false,
	// 				"tags": []string{"tag:prod"},
	// 			},
	// 		},
	// 	},
	// 	"expirySeconds": 86400,
	// 	"description": "short description of key purpose",
	// }	

	// jsonData, formatErr := json.Marshal(data)
	// if formatErr != nil {
	// 	log.Fatalf("error creating body: %v", formatErr)
	// }
	
	// reqBody := bytes.NewBuffer(jsonData)

	// createKeysResp, createKeysErr := client.Post("https://api.tailscale.com/api/v2/tailnet/vungle.com/keys", "application/json", reqBody)
	// if getKeysErr != nil {
	// 	log.Fatalf("error creating keys: %v", createKeysErr)
	// }

	// createKeysBody, createKeysErr := ioutil.ReadAll(createKeysResp.Body)
	// if getKeysErr != nil {
	// 	log.Fatalf("error reading response body: %v", createKeysErr)
	// }
	
	// fmt.Printf("create_keys_body response: %s", string(createKeysBody))		
	
	// Section to call Vault api
	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = "https://vault.ops.vungle.io"

	vaultClient, vaultErr := vault.NewClient(config)
	if vaultErr != nil {
		log.Fatalf("unable to initialize Vault client: %s", vaultErr)
	}

	awsAuth, err := auth.NewAWSAuth(
		auth.WithRole("vault-prod"), // if not provided, Vault will fall back on looking for a role with the IAM role name if you're using the iam auth type, or the EC2 instance's AMI id if using the ec2 auth type
	)
	if err != nil {
		log.Fatalf("unable to initialize AWS auth method: %s", err)
	}
	
	authInfo, err := vaultClient.Auth().Login(context.Background(), awsAuth)
	if err != nil {
		log.Fatalf("unable to login to AWS auth method: %s", err)
	}
	if authInfo == nil {
		log.Fatalf("no auth info was returned after login:")
	}

	// get secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := vaultClient.KVv2("secret").Get(context.Background(), "ops")
	if err != nil {
		log.Fatalf("unable to read secret: %s", err)
	}

	// data map can contain more than one key-value pair,
	// in this case we're just grabbing one of them
	value, ok := secret.Data["test_password"].(string)
	if !ok {
		log.Fatalf("value type assertion failed: %T %#v", secret.Data["test_password"], secret.Data["test_password"])
	}	
	fmt.Printf("value response: %s", value)	
	
	
	// Use token to do
	// vaultClient, err := vault.NewClient(config)
	// if err != nil {
	// 	log.Fatalf("unable to initialize Vault client: %v", err)
	// }	
	
	// // Authenticate
	// vaultClient.SetToken(os.Getenv("PERSONAL_VAULT_TOKEN"))

	// secretData := map[string]interface{}{
	// 	"bingkun_password": "Hashi123",
	// }

	// // Write a secret
	// _, err = vaultClient.KVv2("secret").Put(context.Background(), "bingkun-secret-password", secretData)
	// if err != nil {
	// 	log.Fatalf("unable to write secret: %v", err)
	// }

	// fmt.Println("Secret written successfully.")

	// // Read a secret from the default mount path for KV v2 in dev mode, "secret"
	// secret, err := vaultClient.KVv2("secret").Get(context.Background(), "bingkun-secret-password")
	// if err != nil {
	// 	log.Fatalf("unable to read secret: %v", err)
	// }

	// value, ok := secret.Data["bingkun_password"].(string)
	// if !ok {
	// 	log.Fatalf("value type assertion failed: %T %#v", secret.Data["bingkun_password"], secret.Data["bingkun_password"])
	// }

	// if value != "Hashi123" {
	// 	log.Fatalf("unexpected password value %q retrieved from vault", value)
	// }

	// fmt.Println("Access granted!")	
}

// func getSecretWithAWSAuthIAM() (string, error) {
// 	config := vault.DefaultConfig() // modify for more granular configuration

// 	client, err := vault.NewClient(config)
// 	if err != nil {
// 		return "", fmt.Errorf("unable to initialize Vault client: %w", err)
// 	}

// 	awsAuth, err := auth.NewAWSAuth(
// 		auth.WithRole("vault-prod"), // if not provided, Vault will fall back on looking for a role with the IAM role name if you're using the iam auth type, or the EC2 instance's AMI id if using the ec2 auth type
// 	)
// 	if err != nil {
// 		return "", fmt.Errorf("unable to initialize AWS auth method: %w", err)
// 	}

// 	authInfo, err := client.Auth().Login(context.Background(), awsAuth)
// 	if err != nil {
// 		return "", fmt.Errorf("unable to login to AWS auth method: %w", err)
// 	}
// 	if authInfo == nil {
// 		return "", fmt.Errorf("no auth info was returned after login")
// 	}

// 	// get secret from the default mount path for KV v2 in dev mode, "secret"
// 	secret, err := client.KVv2("secret").Get(context.Background(), "creds")
// 	if err != nil {
// 		return "", fmt.Errorf("unable to read secret: %w", err)
// 	}

// 	// data map can contain more than one key-value pair,
// 	// in this case we're just grabbing one of them
// 	value, ok := secret.Data["test_password"].(string)
// 	if !ok {
// 		return "", fmt.Errorf("value type assertion failed: %T %#v", secret.Data["test_password"], secret.Data["test_password"])
// 	}

// 	return value, nil
// }
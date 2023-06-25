package main

import (
        "context"
        "fmt"
        "io/ioutil"
        "log"
        "os"

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
        resp, err := client.Get("https://api.tailscale.com/api/v2/tailnet/vungle.com/devices")
        if err != nil {
                log.Fatalf("error getting keys: %v", err)
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Fatalf("error reading response body: %v", err)
        }

        fmt.Printf("response: %s", string(body))
}
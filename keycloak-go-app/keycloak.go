package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v11"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type keycloakCreds struct {
	hostname     string
	clientId     string
	clientSecret string
	realm        string
	username     string
	password     string
	grantType    string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var kCreds = &keycloakCreds{
	hostname:     goDotEnvVariable("KEYCLOAK_HOSTNAME"),
	clientId:     goDotEnvVariable("KEYCLOAK_CLIENT_ID"),
	clientSecret: goDotEnvVariable("KEYCLOAK_CLIENT_SECRET"),
	realm:        goDotEnvVariable("KEYCLOAK_REALM"),
	grantType:    goDotEnvVariable("KEYCLOAK_GRANT_TYPE"),
}

func keycloakClientLogin(username string, password string) (string, string, error) {

	kUrl := "https://keycloak.teletek.net.tr/realms/Teletek/protocol/openid-connect/token"

	data := url.Values{}
	data.Set("client_id", kCreds.clientId)
	data.Set("client_secret", kCreds.clientSecret)
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", kCreds.grantType)

	client := resty.New()
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(data.Encode()).
		Post(kUrl)

	if err != nil {
		return "", "", err
	}

	if resp.StatusCode() != 200 {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(resp.Body(), &tokenResponse)
	if err != nil {
		return "", "", err
	}

	return tokenResponse.AccessToken, tokenResponse.RefreshToken, nil
}

func keycloakRetrospectToken(accessToken string) (bool, error) {
	keycloakClient := gocloak.NewClient(kCreds.hostname)
	restyClient := keycloakClient.RestyClient()
	restyClient.SetDebug(true)

	kCTX := context.Background()

	retrospectToken, err := keycloakClient.RetrospectToken(
		kCTX,
		accessToken,
		kCreds.clientId,
		kCreds.clientSecret,
		kCreds.realm,
	)

	if err != nil {
		log.Error().Msgf("%v", "keycloakClient.RetrospectToken() Invalid or malformed token", err)
		return false, err
	}

	if *retrospectToken.Active {
		log.Info().Msgf("%v", "Token is active")
		return true, nil
	}

	return false, nil
}

func keycloakClientTokenRevoke(accessToken string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	endpoint := kCreds.hostname + "auth/realms/" + kCreds.realm + "/protocol/openid-connect/revoke"

	data := url.Values{}
	data.Set("client_id", kCreds.clientId)
	data.Set("client_secret", kCreds.clientSecret)
	data.Set("token", accessToken)
	encodedData := data.Encode()
	fmt.Println(encodedData)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(encodedData))
	if err != nil {
		log.Error().Msgf("%v", "Error creating request", err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	response, err := client.Do(req)

	if err != nil {
		log.Error().Msgf("%v", "Error sending request", err)
		return err
	}
	defer response.Body.Close()
	return nil

}

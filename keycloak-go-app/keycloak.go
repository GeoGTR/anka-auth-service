package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v11"
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

var kCreds = &keycloakCreds{
	hostname:     goDotEnvVariable("KEYCLOAK_HOSTNAME"),
	clientId:     goDotEnvVariable("KEYCLOAK_CLIENT_ID"),
	clientSecret: goDotEnvVariable("KEYCLOAK_CLIENT_SECRET"),
	realm:        goDotEnvVariable("KEYCLOAK_REALM"),
	grantType:    goDotEnvVariable("KEYCLOAK_GRANT_TYPE"),
}

func keycloakClientLogin(username string, password string) (string, string, error) {

	var keycloakClientLoginCreds = &keycloakCreds{
		username: username,
		password: password,
	}

	keycloakClient := gocloak.NewClient(kCreds.hostname)
	restyClient := keycloakClient.RestyClient()
	restyClient.SetDebug(true)

	kCTX := context.Background()

	jwt, err := keycloakClient.GetToken(
		kCTX,
		kCreds.realm,
		gocloak.TokenOptions{
			ClientID:     &kCreds.clientId,
			ClientSecret: &kCreds.clientSecret,
			Username:     &keycloakClientLoginCreds.username,
			Password:     &keycloakClientLoginCreds.password,
			GrantType:    &kCreds.grantType,
		},
	)

	if err != nil {
		log.Error().Msgf("%v", "keycloakClient.Login() Invalid credentials", err)
		return "", "", err
	}

	return jwt.AccessToken, jwt.RefreshToken, nil
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

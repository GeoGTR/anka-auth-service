package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func login(c *gin.Context) {

	var requestBody loginRequest
	var responseBody loginResponse

	if c.FullPath() == "/loginWeb" {
		requestBody.Username = goDotEnvVariable("KEYCLOAK_USERNAME")
		requestBody.Password = goDotEnvVariable("KEYCLOAK_PASSWORD")
	} else {
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	AccessToken, RefreshToken, err := keycloakClientLogin(requestBody.Username, requestBody.Password)

	if err != nil {
		c.JSON(500, gin.H{
			"Status": "Invalid Credentials",
			"Error":  err,
		})
		return
	}

	responseBody.AccessToken = AccessToken
	responseBody.RefreshToken = RefreshToken

	c.JSON(200, gin.H{
		"Status":       "Login Successful",
		"AcessToken":   responseBody.AccessToken,
		"RefreshToken": responseBody.RefreshToken,
	})

}

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"Health": "OK",
	})
}

func status(c *gin.Context) {
	c.JSON(200, gin.H{
		"Status": "OK",
	})
}

func getQuote(c *gin.Context) {
	response, err := http.Get("https://getmeaquote.designedbyaturtle.com/")

	if err != nil {
		fmt.Println(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(responseData))
	c.JSON(200, gin.H{
		"quote": string(responseData),
	})
}

func logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"Status": "Logout Successful",
	})
}

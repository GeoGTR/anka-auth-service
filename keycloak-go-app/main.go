package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var router *gin.Engine

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Error().Msg("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	router = gin.Default()

	initializeRoutes()

	router.Run("localhost:8081")
}

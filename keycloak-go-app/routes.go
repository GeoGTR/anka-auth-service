package main

func initializeRoutes() {

	// unauthenticated routes
	router.GET("/login", login)
	router.GET("/health", health)

	//login alternative
	router.GET("/loginWeb", login)

	// authenticated routes
	authRoute := router.Group("/auth")
	{
		authRoute.GET("/quote", ValidateToken(), getQuote)
		authRoute.GET("/logout", ValidateToken(), TokenRevoke(), logout)
		authRoute.GET("/status", ValidateToken(), status)
	}
}

/*
***

# Application will start from this package

***
*/
package main

import (
	"TestP/pkg/server"
	"TestP/pkg/utilities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var ENV server.Evn

func main() {

	ENV.Port = os.Getenv("PORT")
	if ENV.Port == "" {
		ENV.Port = "9000"
	}
	router := gin.Default()

	// attaching middleware to router
	router.Use(authMiddleware())

	router.POST("/find", ENV.FindIntersection)

	err := router.Run(":" + ENV.Port) //starting server
	if err != nil {
		log.Println("Error starting Server At ", ENV.Port)
		return
	}

}

// authMiddleware is a middleware to parse and validate the token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the JWE token from the request header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		if utilities.Token != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Continue to the next middleware or API handler
		c.Next()
	}
}

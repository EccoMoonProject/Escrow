package main

import (
	"context"
	"escrow/auth"
	"escrow/dispute"
	"escrow/escrow"
	"escrow/moneyStripe"
	"escrow/userData"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set up a connection to the MongoDB database
	clientOptions := options.Client().ApplyURI("mongodb+srv://EscrowBackend:Atlantica123@cluster0.p2vyode.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Create a new Gin router
	router := gin.Default()

	router.Use(CORSMiddleware())

	auth.SetupAuthRoutes(router, client)

	userData.SetupUserRoutes(router, client)

	escrow.SetupEscrowRoutes(router, client)
	escrow.SetupChatsRoutes(router, client)

	moneyStripe.SetupStripeRoutes(router, client)
	moneyStripe.SetupPayoutsRoutes(router, client)

	dispute.SetupDisputeRoutes(router, client)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

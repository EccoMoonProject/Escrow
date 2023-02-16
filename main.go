package main

import (
	"context"
	"escrow/auth"
	"escrow/dispute"
	"escrow/escrow"
	"escrow/moneyStripe"
	"escrow/userData"
	"log"

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

	auth.SetupAuthRoutes(router, client)

	userData.SetupUserRoutes(router, client)

	escrow.SetupEscrowRoutes(router, client)

	moneyStripe.SetupStripeRoutes(router, client)

	dispute.SetupDisputeRoutes(router, client)

	// Start the server
	router.Run(":8080")
}

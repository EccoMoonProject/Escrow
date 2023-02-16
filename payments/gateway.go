package payments

import (
	"context"
	"escrow/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Path: payments/gateway.go

func SetupGatewayRoutes(r *gin.Engine, client *mongo.Client) {

	// Define a GET endpoint to insert data into MongoDB
	r.GET("payments/deposit", func(c *gin.Context) {
		// Get the query parameters
		ownerID := c.Query("ownerID")
		deposit := c.Query("deposit")

		// Convert the deposit to a uint64
		depositUint64, err := strconv.ParseUint(deposit, 10, 64)

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("payments")

		// Create a new Wallet struct with the query parameters
		deposity := types.DepositRequest{OwnerID: ownerID, Amount: depositUint64}

		// Insert the wallet into the wallets collection
		_, err = collection.InsertOne(context.Background(), deposity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert deposit"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deposit inserted successfully"})

	})

	r.GET("payments/withdraw", func(c *gin.Context) {
		// Get the query parameters
		ownerID := c.Query("ownerID")
		deposit := c.Query("deposit")

		// Convert the deposit to a uint64
		depositUint64, err := strconv.ParseUint(deposit, 10, 64)

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("payments")

		// Create a new Wallet struct with the query parameters
		withdraw := types.WithdrawRequest{OwnerID: ownerID, Amount: depositUint64}

		// Insert the wallet into the wallets collection
		_, err = collection.InsertOne(context.Background(), withdraw)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert deposit"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deposit inserted successfully"})
	})

}

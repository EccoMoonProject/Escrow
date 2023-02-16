package userData

import (
	"context"
	"escrow/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupUserRoutes(r *gin.Engine, client *mongo.Client) {

	r.GET("user/wallet/createWallet", func(c *gin.Context) {
		// Get the query parameters
		ownerID := c.Query("ownerID")
		deposit := c.Query("deposit")

		// Convert the deposit to a uint64
		depositUint64, err := strconv.ParseUint(deposit, 10, 64)

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a new Wallet struct with the query parameters
		wallet := types.Wallet{OwnerID: ownerID, Balance: depositUint64}

		// Insert the wallet into the wallets collection
		_, err = collection.InsertOne(context.Background(), wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet inserted successfully"})
	})

	r.GET("user/wallet/getBalance", func(c *gin.Context) {
		ownerID := c.Query("ownerID")

		// Get a handle to the wallets collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": ownerID}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err := collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"balance": wallet.Balance})
	})

	r.GET("user/wallet/deposit", func(c *gin.Context) {
		// Get the query parameters
		ownerID := c.Query("ownerID")
		deposit := c.Query("deposit")

		// Convert the deposit to a uint64
		depositUint64, err := strconv.ParseUint(deposit, 10, 64)

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": ownerID}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err = collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		// update the wallet balance
		update := bson.M{"$set": bson.M{"balance": wallet.Balance + depositUint64}}

		// Update the wallet balance
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})
	})

	r.GET("user/wallet/withdraw", func(c *gin.Context) {
		// Get the query parameters
		ownerID := c.Query("ownerID")
		withdrawal := c.Query("withdrawal")

		// Convert the withdrawal to a uint64
		withdrawalUint64, err := strconv.ParseUint(withdrawal, 10, 64)

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": ownerID}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err = collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		// Check if the wallet has enough balance
		if wallet.Balance < withdrawalUint64 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insufficient balance"})
			return
		}

		// update the wallet balance
		update := bson.M{"$set": bson.M{"balance": wallet.Balance - withdrawalUint64}}

		// Update the wallet balance
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})
	})

	r.GET("user/wallet/paymentRequest", func(c *gin.Context) {
		ownerID := c.Query("ownerID")
		amount := c.Query("amount")

		// Convert the amount to a uint64
		amountUint64, err := strconv.ParseUint(amount, 10, 64)

		// Get a handle to the wallets collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": ownerID}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err = collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		// Check if the wallet has enough balance
		if wallet.Balance < amountUint64 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insufficient balance"})
			return
		}

		// update the wallet balance
		update := bson.M{"$set": bson.M{"balance": wallet.Balance - amountUint64}}

		// Update the wallet balance
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})

	})

}

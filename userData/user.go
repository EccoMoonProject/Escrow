package userData

import (
	"context"
	"escrow/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupUserRoutes(r *gin.Engine, client *mongo.Client) {

	r.POST("user/wallet/createWallet", func(c *gin.Context) {

		var walletReq types.Wallet
		err := c.BindJSON(&walletReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a new Wallet struct with the query parameters
		wallet := types.Wallet{OwnerID: walletReq.OwnerID, Balance: walletReq.Balance}

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

	r.POST("user/wallet/deposit", func(c *gin.Context) {

		var walletReq types.Wallet
		err := c.BindJSON(&walletReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": walletReq.OwnerID}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err = collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		// update the wallet balance
		update := bson.M{"$set": bson.M{"balance": wallet.Balance + walletReq.Balance}}

		// Update the wallet balance
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})
	})

	r.POST("user/wallet/withdraw", func(c *gin.Context) {

		var walletReq types.Wallet
		err := c.BindJSON(&walletReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": walletReq.OwnerID}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err = collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		// Check if the wallet has enough balance
		if wallet.Balance < walletReq.Balance {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insufficient balance"})
			return
		}

		// update the wallet balance
		update := bson.M{"$set": bson.M{"balance": wallet.Balance - walletReq.Balance}}

		// Update the wallet balance
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})
	})

	r.POST("user/wallet/paymentRequest", func(c *gin.Context) {

		var walletReq types.Wallet
		err := c.BindJSON(&walletReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		// Get a handle to the wallets collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": walletReq.OwnerID}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err = collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		// Check if the wallet has enough balance
		if wallet.Balance < walletReq.Balance {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insufficient balance"})
			return
		}

		// update the wallet balance
		update := bson.M{"$set": bson.M{"balance": wallet.Balance - walletReq.Balance}}

		// Update the wallet balance
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet updated successfully"})

	})

	r.GET("user/wallet/getWallet", func(c *gin.Context) {
		owner_id := c.Query("ownerID")

		// Get a handle to the wallets collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": owner_id}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err := collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find wallet"})
			return
		}

		c.JSON(http.StatusOK, wallet)
	})

	r.GET("user/wallet/isWalletCreated", func(c *gin.Context) {
		owner_id := c.Query("ownerID")

		// Get a handle to the wallets collection
		collection := client.Database("mydb").Collection("wallets")

		// Create a filter to find the wallet with the ownerID
		filter := bson.M{"ownerID": owner_id}

		// Find the wallet with the ownerID
		var wallet types.Wallet
		err := collection.FindOne(context.Background(), filter).Decode(&wallet)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"isWalletCreated": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"isWalletCreated": true})
	})

	r.GET("user/getUser", func(c *gin.Context) {
		owner_id := c.Query("ownerID")

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("users")

		// Create a filter to find the user with the userID
		filter := bson.M{"ownerID": owner_id}

		// Find the user with the userID
		var user types.User
		err := collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
			return
		}

		c.JSON(http.StatusOK, user)
	})
}

package dispute

import (
	"context"
	"escrow/types"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupDisputeRoutes(r *gin.Engine, client *mongo.Client) {

	r.GET("dispute/createDispute", func(c *gin.Context) {
		ownerID := c.Query("owner_id")
		ownerSHI := c.Query("owner_shi")
		category := c.Query("category")
		amount := c.Query("amount")

		// Convert the amount to a uint64
		amountUint64, err := strconv.ParseUint(amount, 10, 64)

		dispute := types.Dispute{OwnerID: ownerID, OwnerSHI: ownerSHI, Category: category, Amount: amountUint64}

		// Get a handle to the escrowInstances collection
		collection := client.Database("mydb").Collection("dispute")

		// Insert the escrowInstance into the escrowInstances collection
		_, err = collection.InsertOne(context.Background(), dispute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert dispute"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Dispute inserted successfully"})

	})
	// Endpoint to create a chat
	r.POST("dispute/chats", func(c *gin.Context) {
		// Parse the request body into a Chat struct
		var chat types.Chat
		err := c.BindJSON(&chat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		// Generate a 5-letter random string
		letters := "abcdefghijklmnopqrstuvwxyz"
		result := make([]byte, 5)
		for i := range result {
			result[i] = letters[rand.Intn(len(letters))]
		}
		randomString := fmt.Sprintf("%s", result)

		chat.ID = randomString

		// Get a handle to the chats collection
		collection := client.Database("mydb").Collection("chats")

		// Insert the chat into the chats collection
		_, err = collection.InsertOne(context.Background(), chat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert chat"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully"})
	})

	r.GET("dispute/getChats", func(c *gin.Context) {
		// Get the sender and ID query parameters
		sender := c.Query("sender")
		id := c.Query("owner_shi")

		// Get a handle to the chats collection
		collection := client.Database("mydb").Collection("chats")

		// Define a filter to retrieve chat messages by sender and id
		filter := bson.M{
			"sender":   sender,
			"ownerSHI": id,
		}

		// Retrieve chat messages by sender and id
		cur, err := collection.Find(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chat history"})
			return
		}

		// Build a slice of chat messages
		var chats []types.Chat
		for cur.Next(context.Background()) {
			var chat types.Chat
			err := cur.Decode(&chat)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode chat message"})
				return
			}
			chats = append(chats, chat)
		}

		// Return the chat messages in JSON format
		c.JSON(http.StatusOK, chats)
	})
}

package escrow

import (
	"context"
	"crypto/sha256"
	"escrow/types"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupEscrowRoutes(r *gin.Engine, client *mongo.Client) {

	r.POST("escrow/createInstance", func(c *gin.Context) {

		var escrowRequest types.EscrowRequest
		err := c.BindJSON(&escrowRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		// Initialize the random number generator
		rand.Seed(time.Now().UnixNano())

		// Generate a 5-letter random string
		letters := "abcdefghijklmnopqrstuvwxyz"
		result := make([]byte, 5)
		for i := range result {
			result[i] = letters[rand.Intn(len(letters))]
		}
		randomString := fmt.Sprintf("%s", result)

		fmt.Println(randomString) // Prints a random 5-letter string

		// Compute the SHA256 hash of the string
		hash := sha256.Sum256([]byte(randomString))

		// Convert the hash to a string
		secure_hash_identifier := fmt.Sprintf("%x", hash)

		fmt.Println(secure_hash_identifier) // Prints the SHA256 hash of the random string
		var secureDestroyer bool = false
		var status bool = false

		// Create a new EscrowInstance struct with the query parameters
		escrowInstance := types.EscrowInstance{InstanceID: randomString, OwnerID: escrowRequest.OwnerID, OwnerEmail: escrowRequest.Email, Amount: escrowRequest.Amount, Status: status, OwnerSHI: secure_hash_identifier, SecureDestroyer: secureDestroyer}

		// Get a handle to the escrowInstances collection
		collection := client.Database("mydb").Collection("escrowInstances")

		// Insert the escrowInstance into the escrowInstances collection
		_, err = collection.InsertOne(context.Background(), escrowInstance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert escrowInstance"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Escrow Instance created", "instanceId": randomString, "ownerId": escrowRequest.OwnerID, "ownerSHI": secure_hash_identifier, "currency": escrowRequest.Currency, "amount": escrowRequest.Amount, "status": status, "secureDestroyer": secureDestroyer})

	})

	r.GET("escrow/destroyInstance", func(c *gin.Context) {
		// Get the query parameters
		shi := c.Query("ownerSHI")

		// Get a handle to the escrowInstances collection
		collection := client.Database("mydb").Collection("escrowInstances")
		// Delete the user by email
		filter := bson.M{"ownerSHI": shi}
		result, err := collection.DeleteOne(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting instance"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Instance deleted"})
	})

	r.POST("escrow/payInstance", func(c *gin.Context) {

		var p types.PayInstanceRequest
		err := c.BindJSON(&p)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		// find the instance
		collection := client.Database("mydb").Collection("escrowInstances")

		// Create a filter to find the instance
		filter := bson.M{"instanceID": p.InstanceID, "ownerID": p.OwnerID, "ownerSHI": p.OwnerSHI}

		// Find the instance
		var escrowInstance types.EscrowInstance
		err = collection.FindOne(context.Background(), filter).Decode(&escrowInstance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding instance"})
			return
		}
		// get the amount of the instance
		amountInstance := escrowInstance.Amount

		// check if the amount is enough
		if p.Amount < amountInstance {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not enough amount"})
			return
		}

		// update the instance
		update := bson.M{"$set": bson.M{"status": true, "secureDestroyer": true}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating instance"})
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Instance updated"})

	})

	r.GET("escrow/getInstanceStatus", func(c *gin.Context) {
		// Get the query parameters
		instanceId := c.Query("instanceId")

		// Get a handle to the escrowInstances collection
		collection := client.Database("mydb").Collection("escrowInstances")
		// Delete the user by email
		filter := bson.M{"instanceID": instanceId}
		var escrowInstance types.EscrowInstance
		err := collection.FindOne(context.Background(), filter).Decode(&escrowInstance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding instance"})
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Instance found", "status": escrowInstance.Status})
	})

	r.POST("escrow/voting/createPool", func(c *gin.Context) {

		var p types.VotingRequest
		var buyerVote uint16 = 0
		var sellerVote uint16 = 0
		var consensus bool = false

		// Create a new VotingPool struct with the query parameters
		votingPool := types.VotingPool{InstanceID: p.InstanceID, BuyerID: p.BuyerID, SellerID: p.SellerID, BuyerVote: buyerVote, SellerVote: sellerVote, Consensus: consensus}

		// Get a handle to the votingPools collection
		collection := client.Database("mydb").Collection("votingPools")

		// Insert the votingPool into the votingPools collection
		_, err := collection.InsertOne(context.Background(), votingPool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert votingPool"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "VotingPool inserted successfully"})
	})

	r.GET("escrow/voting/vote", func(c *gin.Context) {
		instanceId := c.Query("instanceId")
		userId := c.Query("userId")
		vote := c.Query("vote")

		// Convert the vote to a uint16
		voteUint16, err := strconv.ParseUint(vote, 10, 16)

		// find the votingPool
		collection := client.Database("mydb").Collection("votingPools")

		// Create a filter to find the votingPool
		filter := bson.M{"instanceID": instanceId}

		// Find the votingPool
		var votingPool types.VotingPool
		err = collection.FindOne(context.Background(), filter).Decode(&votingPool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding votingPool"})
			return
		}

		// check if the user is the buyer or the seller
		if votingPool.BuyerID == userId {
			// update the votingPool
			update := bson.M{"$set": bson.M{"buyerVote": voteUint16}}
			_, err = collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating votingPool"})
				return
			}
		} else if votingPool.SellerID == userId {
			// update the votingPool
			update := bson.M{"$set": bson.M{"sellerVote": voteUint16}}
			_, err = collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating votingPool"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User is not the buyer or the seller"})
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "VotingPool updated"})
	})

	r.GET("escrow/voting/consensus", func(c *gin.Context) {
		// if the buyer and the seller voted the same, consensus is true
		instanceId := c.Query("instanceId")

		// find the votingPool
		collection := client.Database("mydb").Collection("votingPools")

		// Create a filter to find the votingPool
		filter := bson.M{"instanceID": instanceId}

		// Find the votingPool
		var votingPool types.VotingPool

		err := collection.FindOne(context.Background(), filter).Decode(&votingPool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding votingPool"})
			return
		}

		// check if the buyer and the seller voted the same
		if votingPool.BuyerVote == votingPool.SellerVote {
			// update the votingPool
			update := bson.M{"$set": bson.M{"consensus": true}}
			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating votingPool"})
				return
			}
		} else {
			// update the votingPool
			update := bson.M{"$set": bson.M{"consensus": false}}
			_, err := collection.UpdateOne(context.Background(), filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating votingPool"})
				return
			}
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "VotingPool updated"})
	})

	r.GET("escrow/voting/getConsensus", func(c *gin.Context) {
		instanceId := c.Query("instanceId")

		// find the votingPool
		collection := client.Database("mydb").Collection("votingPools")

		// Create a filter to find the votingPool
		filter := bson.M{"instanceID": instanceId}

		// Find the votingPool
		var votingPool types.VotingPool
		err := collection.FindOne(context.Background(), filter).Decode(&votingPool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding votingPool"})
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "VotingPool found", "consensus": votingPool.Consensus})
	})
}

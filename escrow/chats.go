package escrow

import (
	"context"
	"escrow/types"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupChatsRoutes(r *gin.Engine, client *mongo.Client) {

	r.POST("chats/createChat", func(c *gin.Context) {
		var chatRoom types.ChatRoom
		if err := c.ShouldBindJSON(&chatRoom); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		chatRoom.ID = primitive.NewObjectID()
		chatRoomsCollection := client.Database("mydb").Collection("chat_rooms")
		if _, err := chatRoomsCollection.InsertOne(context.Background(), chatRoom); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": chatRoom.ID})

	})

	r.POST("chats/:id/join", func(c *gin.Context) {
		chatRoomID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat room ID"})
			return
		}
		var joinRequest struct {
			Username string `json:"username"`
		}
		if err := c.ShouldBindJSON(&joinRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"_id": chatRoomID}
		update := bson.M{"$addToSet": bson.M{"members": joinRequest.Username}}
		chatRoomsCollection := client.Database("mydb").Collection("chat_rooms")
		if _, err := chatRoomsCollection.UpdateOne(context.Background(), filter, update); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Joined chat room"})
	})

	// Define the Gin route for sending a message in a chat room
	r.POST("/chatrooms/:id/messages", func(c *gin.Context) {
		chatRoomID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat room ID"})
			return
		}
		var message struct {
			Username   string             `json:"username"`
			Content    string             `json:"content"`
			Timestamp  primitive.DateTime `json:"timestamp"`
			ChatRoomID primitive.ObjectID `json:"chatRoomId"`
		}
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{"_id": chatRoomID}
		update := bson.M{"$push": bson.M{"messages": message}}
		chatRoomsCollection := client.Database("mydb").Collection("chat_rooms")
		if _, err := chatRoomsCollection.UpdateOne(context.Background(), filter, update); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		message.Timestamp = primitive.NewDateTimeFromTime(time.Now().UTC())
		message.ChatRoomID = chatRoomID
		messagesCollection := client.Database("mydb").Collection("messages")
		if _, err := messagesCollection.InsertOne(context.Background(), message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Sent message"})
	})

	// Define the Gin route for joining a chat room
	r.GET("/chatrooms/:id", func(c *gin.Context) {
		chatRoomID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat room ID"})
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Create a WaitGroup to keep track of the number of clients listening to the chat room
		var wg sync.WaitGroup
		wg.Add(1)

		// Create a channel to receive messages from the chat room
		messageChan := make(chan interface{})
		defer close(messageChan)

		// Start a goroutine to listen for new messages in the chat room
		go func() {
			defer wg.Done()

			// Create a MongoDB change stream to watch for new messages in the chat room
			pipeline := mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"operationType": "insert", "fullDocument.chatRoomID": chatRoomID.Hex()}}},
			}
			streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)
			messagesCollection := client.Database("mydb").Collection("messages")
			changeStream, err := messagesCollection.Watch(ctx, pipeline, streamOptions)
			if err != nil {
				// Handle error
			}
			defer changeStream.Close(ctx)

			// Loop through the change stream and send messages to the message channel
			for changeStream.Next(ctx) {
				var message bson.M
				if err := changeStream.Decode(&message); err != nil {
					// Handle error
				}
				message["_id"] = message["_id"].(primitive.ObjectID).Hex()
				message["timestamp"] = message["timestamp"].(primitive.DateTime).Time().UTC().Format(time.RFC3339Nano)
				messageChan <- message
			}
			if err := changeStream.Err(); err != nil {
				// Handle error
			}
		}()

		// Start a goroutine to send messages to the client
		go func() {
			// Set a timer for 5 minutes
			timeout := time.NewTimer(time.Minute * 5)
			defer timeout.Stop()

			for {
				select {
				case message := <-messageChan:
					c.SSEvent("message", message)
					// Reset the timer for another 5 minutes
					timeout.Reset(time.Minute * 5)
				case <-timeout.C:
					// Disconnect the client if there is no message received from the message channel within 5 minutes
					c.AbortWithStatus(http.StatusNoContent)
					return
				}
			}
		}()

		// Wait for the client to disconnect before canceling the context
		c.Writer.WriteHeader(http.StatusOK)
		wg.Wait()
	})

	// Define the Gin route for retrieving chat room messages
	r.GET("/chatrooms/:id/messages", func(c *gin.Context) {
		chatRoomID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat room ID"})
			return
		}

		// Define a filter to retrieve messages for the specified chat room
		filter := bson.M{"chatroomid": chatRoomID}

		// Retrieve messages from the messages collection
		messagesCollection := client.Database("mydb").Collection("messages")
		cursor, err := messagesCollection.Find(context.Background(), filter)
		if err != nil {
			// Handle error
		}
		defer cursor.Close(context.Background())
		// Loop through the cursor and add each message to the response array
		var messages []bson.M
		for cursor.Next(context.Background()) {
			var message bson.M
			if err := cursor.Decode(&message); err != nil {
				// Handle error
			}
			message["_id"] = message["_id"].(primitive.ObjectID).Hex()
			message["timestamp"] = message["timestamp"].(primitive.DateTime).Time().UTC().Format(time.RFC3339Nano)
			messages = append(messages, message)
		}
		if err := cursor.Err(); err != nil {
			// Handle error
		}

		// Return the messages array in the response
		c.JSON(http.StatusOK, messages)
	})

}

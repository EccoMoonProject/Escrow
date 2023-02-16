package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name        string `bson:"name"`
	Email       string `bson:"email"`
	Password    string `bson:"password"`
	DateOfBirth string `bson:"dateOfBirth"`
}

func SetupAuthRoutes(r *gin.Engine, client *mongo.Client) {

	// Define a GET endpoint to insert data into MongoDB
	r.GET("auth/signup", func(c *gin.Context) {
		// Get the query parameters
		name := c.Query("name")
		email := c.Query("email")
		password := c.Query("password")
		dob := c.Query("dateOfBirth")

		// Create a new User struct with the query parameters
		user := User{Name: name, Email: email, Password: password, DateOfBirth: dob}

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("users")

		// Insert the user into the users collection
		_, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User inserted successfully"})
	})

	r.GET("auth/login", func(c *gin.Context) {

		// Get the query parameters
		email := c.Query("email")
		password := c.Query("password")

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("users")

		// Create a filter to find the user with the email
		filter := bson.M{"email": email}

		// Find the user with the email
		var user User
		err := collection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
			return
		}

		// Check if the password matches
		if password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
	})

}

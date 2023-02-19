package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name     string `bson:"name"`
	OwnerID  string `bson:"ownerID"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type Signup struct {
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type Login struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

func SetupAuthRoutes(r *gin.Engine, client *mongo.Client) {

	// Define a GET endpoint to insert data into MongoDB
	r.POST("auth/signup", func(c *gin.Context) {

		var signup Signup
		err := c.BindJSON(&signup)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		fmt.Println(signup)

		// generate a random ownerID
		letters := "abcdefghijklmnopqrstuvwxyz"
		result := make([]byte, 8)
		for i := range result {
			result[i] = letters[rand.Intn(len(letters))]
		}
		randomString := fmt.Sprintf("%s", result)

		// make sha256 and add salt
		hashed := sha256.Sum256([]byte(signup.Password + randomString))

		// convert to string
		owner_id := fmt.Sprintf("%x", hashed)

		// reduce to length 5 letters
		owner_id = owner_id[:7]

		// make sha256 from password
		hashed = sha256.Sum256([]byte(signup.Password))

		// convert to string

		signup.Password = fmt.Sprintf("%x", hashed)

		// Create a new User struct with the query parameters
		user := User{Name: signup.Name, OwnerID: owner_id, Email: signup.Email, Password: signup.Password}

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("users")

		// check if the user already exists
		filter := bson.M{"email": user.Email}
		var userExists User
		errorl := collection.FindOne(context.Background(), filter).Decode(&userExists)
		if errorl == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists"})
			return
		}

		// Insert the user into the users collection
		_, errorly := collection.InsertOne(context.Background(), user)
		if errorly != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": owner_id, "state": true})
	})

	r.POST("auth/login", func(c *gin.Context) {

		var login Login
		err := c.BindJSON(&login)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		fmt.Println(login)

		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("users")

		// Create a filter to find the user with the email
		filter := bson.M{"email": login.Email}

		// Find the user with the email
		var user User
		errorly := collection.FindOne(context.Background(), filter).Decode(&user)
		if errorly != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
			return
		}

		// Check if the password matches but first make sha256 from password
		hashed := sha256.Sum256([]byte(login.Password))

		// convert to string

		login.Password = fmt.Sprintf("%x", hashed)
		if login.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": user.OwnerID, "state": true})
	})

}

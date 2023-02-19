package moneyStripe

import (
	"escrow/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/bankaccount"
	"github.com/stripe/stripe-go/v74/card"
	"github.com/stripe/stripe-go/v74/payout"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupPayoutsRoutes(r *gin.Engine, client *mongo.Client) {
	stripe.Key = "sk_test_51Mc3KYI6g21tTAothF9T3fa6LKoGMogNxkYda1qvj1NjrAbZ1FtrNOticsTNNpmFZStHzEW4WjTo6zsNK7yeWpIS00qCnTyzfv"

	r.POST("payouts/createBankAccount", func(c *gin.Context) {
		var b types.BankAccount
		err := c.BindJSON(&b)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}
		params := &stripe.BankAccountParams{
			Customer:          stripe.String(b.Customer),
			AccountHolderName: stripe.String(b.AccountHolderName),
			AccountHolderType: stripe.String(string(stripe.BankAccountAccountHolderTypeIndividual)),
			Country:           stripe.String(string(b.Country)),
			Currency:          stripe.String(string(b.Currency)),
			AccountNumber:     stripe.String(b.AccountNumber),
		}
		ba, err := bankaccount.New(params)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// insert bank account into database
		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("bank_accounts")

		_, err = collection.InsertOne(c, ba)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert bank account"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"bank_account": ba})

	})

	r.POST("payouts/createCard", func(c *gin.Context) {
		var b types.CreateCardRequest
		err := c.BindJSON(&b)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}
		params := &stripe.CardParams{
			Customer: stripe.String(b.Customer),
			// Name:     stripe.String(b.Name),
			// Currency: stripe.String(string(b.Currency)),
			// Number:   stripe.String(b.Number),
			// ExpMonth: stripe.String(b.ExpMonth),
			// ExpYear:  stripe.String(b.ExpYear),
			// CVC:      stripe.String(b.CVC),
			Token: stripe.String("tok_mastercard"),
		}
		card, err := card.New(params)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// insert card into database
		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("cards")

		// create object to insert into database
		cardToInsert := types.CardData{Card: card, Customer: b.Customer}

		_, err = collection.InsertOne(c, cardToInsert)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert card"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"card": card})

	})

	r.GET("payouts/getCardByCustomer", func(c *gin.Context) {
		customer := c.Query("customer")
		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("cards")

		// Create a filter to find the user with the email
		filter := bson.M{"customer": customer}

		var card types.CardData
		err := collection.FindOne(c, filter).Decode(&card)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find card"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"card": card})

	})

	r.POST("payouts/createPayout", func(c *gin.Context) {
		var pyt types.Payouts
		err := c.BindJSON(&pyt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}
		pyt.Amount = pyt.Amount * 100
		params := &stripe.PayoutParams{
			Amount:      stripe.Int64(pyt.Amount),
			Currency:    stripe.String(string(pyt.Currency)),
			Destination: stripe.String(pyt.Destination),
			Method:      stripe.String(string(stripe.PayoutMethodInstant)),
		}
		p, _ := payout.New(params)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// insert payout into database
		// Get a handle to the users collection
		collection := client.Database("mydb").Collection("payouts")

		_, err = collection.InsertOne(c, p)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert payout"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"payout": p})
	})

}

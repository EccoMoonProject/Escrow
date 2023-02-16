package moneyStripe

import (
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/balance"
	"github.com/stripe/stripe-go/balancetransaction"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupStripeRoutes(r *gin.Engine, client *mongo.Client) {

	stripe.Key = "sk_test_51Mc3KYI6g21tTAothF9T3fa6LKoGMogNxkYda1qvj1NjrAbZ1FtrNOticsTNNpmFZStHzEW4WjTo6zsNK7yeWpIS00qCnTyzfv"

	b, _ := balance.Get(nil)
	bt, _ := balancetransaction.Get(
		"txn_3Mc3ZPI6g21tTAot1xF24Dj9",
		nil,
	)

	r.GET("money/getBalance", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"balance": b,
		})
	})

	r.GET("money/getBalanceTransaction", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"balanceTransaction": bt,
		})
	})
}

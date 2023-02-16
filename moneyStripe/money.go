package moneyStripe

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentlink"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupStripeRoutes(r *gin.Engine, client *mongo.Client) {

	stripe.Key = "sk_test_51Mc3KYI6g21tTAothF9T3fa6LKoGMogNxkYda1qvj1NjrAbZ1FtrNOticsTNNpmFZStHzEW4WjTo6zsNK7yeWpIS00qCnTyzfv"

	r.GET("money/createProduct", func(c *gin.Context) {
		productName := c.Query("product")
		params := &stripe.ProductParams{
			Name: stripe.String(productName),
		}
		p, _ := product.New(params)

		c.JSON(200, gin.H{
			"product": p,
		})
	})

	r.GET("money/createPrice", func(c *gin.Context) {
		productId := c.Query("product")
		amount := c.Query("amount")
		currency := c.Query("currency")

		// convert amount to int64
		amountInt64, _ := strconv.ParseInt(amount, 10, 64)

		params := &stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			Product:    stripe.String(productId),
			UnitAmount: stripe.Int64(amountInt64),
		}
		p, _ := price.New(params)

		c.JSON(200, gin.H{
			"price": p,
		})

	})

	r.GET("money/createPaymentLink", func(c *gin.Context) {
		priceId := c.Query("price")
		params := &stripe.PaymentLinkParams{
			LineItems: []*stripe.PaymentLinkLineItemParams{
				{
					Price:    stripe.String(priceId),
					Quantity: stripe.Int64(1),
				},
			},
		}
		pl, _ := paymentlink.New(params)

		c.JSON(200, gin.H{
			"paymentLink": pl,
		})
	})

	r.GET("money/createPaymentIntent", func(c *gin.Context) {
		amount := c.Query("amount")
		destination := c.Query("destination")

		// convert amount to int64
		amountInt64, _ := strconv.ParseInt(amount, 10, 64)

		params := &stripe.PaymentIntentParams{
			Amount: stripe.Int64(amountInt64),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
			Currency: stripe.String(string(stripe.CurrencyPLN)),
			TransferData: &stripe.PaymentIntentTransferDataParams{
				Destination: stripe.String(destination),
			},
		}
		pi, _ := paymentintent.New(params)

		c.JSON(200, gin.H{
			"paymentIntent": pi,
		})
	})

	r.GET("money/payIntent", func(c *gin.Context) {
		paymentId := c.Query("paymentId")

		// To create a PaymentIntent for confirmation, see our guide at: https://stripe.com/docs/payments/payment-intents/creating-payment-intents#creating-for-automatic
		params := &stripe.PaymentIntentConfirmParams{
			PaymentMethod: stripe.String("pm_card_visa"),
		}
		pi, _ := paymentintent.Confirm(
			paymentId,
			params,
		)

		c.JSON(200, gin.H{
			"success": pi,
		})

	})

	r.GET("money/createCheckout", func(c *gin.Context) {
		price := c.Query("price")

		params := &stripe.CheckoutSessionParams{
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(price),
					Quantity: stripe.Int64(1),
				},
			},
			Mode:       stripe.String("payment"),
			SuccessURL: stripe.String("https://example.com/success"),
		}
		s, _ := session.New(params)

		c.JSON(200, gin.H{
			"checkout": s,
		})
	})

}

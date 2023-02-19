package moneyStripe

import (
	"escrow/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/billingportal/session"
	"github.com/stripe/stripe-go/v74/charge"
	sessionCheckout "github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/customerbalancetransaction"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentlink"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
	"github.com/stripe/stripe-go/v74/transfer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupStripeRoutes(r *gin.Engine, client *mongo.Client) {

	stripe.Key = "sk_test_51Mc3KYI6g21tTAothF9T3fa6LKoGMogNxkYda1qvj1NjrAbZ1FtrNOticsTNNpmFZStHzEW4WjTo6zsNK7yeWpIS00qCnTyzfv"

	r.POST("money/createProduct", func(c *gin.Context) {
		var productInstance types.Product
		err := c.BindJSON(&productInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		params := &stripe.ProductParams{
			Name: stripe.String(productInstance.Name),
		}
		p, _ := product.New(params)

		// insert product into db
		collection := client.Database("mydb").Collection("products")
		_, err = collection.InsertOne(c, p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product"})
			return
		}

		c.JSON(200, gin.H{
			"product": p,
		})
	})

	r.POST("money/createPrice", func(c *gin.Context) {

		var priceInstance types.Price
		err := c.BindJSON(&priceInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		priceInstance.Amount = priceInstance.Amount * 100

		params := &stripe.PriceParams{
			Currency:   stripe.String(string(priceInstance.Currency)),
			Product:    stripe.String(priceInstance.Product),
			UnitAmount: stripe.Int64(priceInstance.Amount),
		}
		p, _ := price.New(params)

		// insert price into db
		collection := client.Database("mydb").Collection("prices")

		_, err = collection.InsertOne(c, p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert price"})
			return
		}

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

		// insert payment link into db
		collection := client.Database("mydb").Collection("paymentLinks")

		_, err := collection.InsertOne(c, pl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert payment link"})
			return
		}

		c.JSON(200, gin.H{
			"paymentLink": pl,
		})
	})

	r.GET("money/createPaymentIntent", func(c *gin.Context) {
		amount := c.Query("amount")
		customer := c.Query("customerID")

		// convert amount to int64
		amountInt64, _ := strconv.ParseInt(amount, 10, 64)

		// calculate fee and subtract from amount
		fee := float64(amountInt64) * 0.005
		netAmount := amountInt64 - int64(fee)

		params := &stripe.PaymentIntentParams{
			Amount: stripe.Int64(netAmount),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
			Currency: stripe.String(string(stripe.CurrencyPLN)),
			Customer: stripe.String(customer),
		}
		pi, _ := paymentintent.New(params)

		c.JSON(200, gin.H{
			"paymentIntent": pi,
			"fee":           fee,
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
		s, _ := sessionCheckout.New(params)

		c.JSON(200, gin.H{
			"checkout": s,
		})
	})

	r.POST("money/createCustomer", func(c *gin.Context) {
		var customerInstance types.Customer
		err := c.BindJSON(&customerInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		params := &stripe.CustomerParams{
			Email: stripe.String(customerInstance.Email),
			Name:  stripe.String(customerInstance.OwnerID),
		}
		customer, _ := customer.New(params)

		// insert customer into db
		collection := client.Database("mydb").Collection("customers")

		_, err = collection.InsertOne(c, customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert customer"})
			return
		}

		c.JSON(200, gin.H{
			"customer": customer,
		})

	})

	r.GET("money/getCustomer", func(c *gin.Context) {
		name := c.Query("name")

		collection := client.Database("mydb").Collection("customers")

		var customer *stripe.Customer
		err := collection.FindOne(c, bson.M{"name": name}).Decode(&customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find customer"})
			return
		}

		c.JSON(200, gin.H{
			"customer": customer,
		})

	})

	r.POST("money/createCustomerBalanceTransaction", func(c *gin.Context) {

		var customerBalanceTransactionInstance types.Transfer
		err := c.BindJSON(&customerBalanceTransactionInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		customerBalanceTransactionInstance.Amount = customerBalanceTransactionInstance.Amount * 100

		params := &stripe.
			CustomerBalanceTransactionParams{
			Amount:   stripe.Int64(customerBalanceTransactionInstance.Amount),
			Currency: stripe.String(string(customerBalanceTransactionInstance.Currency)),
			Customer: stripe.String(customerBalanceTransactionInstance.Dest),
		}
		bt, _ := customerbalancetransaction.New(params)

		// insert customer balance transaction into db
		collection := client.Database("mydb").Collection("customerBalanceTransactions")

		_, err = collection.InsertOne(c, bt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert customer balance transaction"})
			return
		}

		c.JSON(200, gin.H{
			"customerBalanceTransaction": bt,
		})
	})

	r.POST("money/createCustomerPortal", func(c *gin.Context) {
		var customerPortalInstance types.CustomerPortal
		err := c.BindJSON(&customerPortalInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		params := &stripe.BillingPortalSessionParams{
			Customer:  stripe.String(customerPortalInstance.Customer),
			ReturnURL: stripe.String("http://localhost:4200/user"),
		}
		s, _ := session.New(params)

		c.JSON(200, gin.H{
			"customerPortal": s,
		})

	})

	r.POST("money/createPaymentMethod", func(c *gin.Context) {
		var paymentMethodInstance types.PaymentMethod
		err := c.BindJSON(&paymentMethodInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		params := &stripe.PaymentMethodParams{
			Card: &stripe.PaymentMethodCardParams{
				Number:   stripe.String(paymentMethodInstance.Card.Number),
				ExpMonth: stripe.Int64(paymentMethodInstance.Card.ExpMonth),
				ExpYear:  stripe.Int64(paymentMethodInstance.Card.ExpYear),
				CVC:      stripe.String(paymentMethodInstance.Card.CVC),
			},
			Type: stripe.String("card"),
		}
		pm, _ := paymentmethod.New(params)

		// add customer to payment method
		params2 := &stripe.PaymentMethodAttachParams{
			Customer: stripe.String(paymentMethodInstance.CustomerID),
		}
		pm, _ = paymentmethod.Attach(
			pm.ID,
			params2,
		)

		// insert payment method into db
		collection := client.Database("mydb").Collection("paymentMethods")

		_, err = collection.InsertOne(c, pm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert payment method"})
			return
		}

		c.JSON(200, gin.H{
			"paymentMethod": pm,
		})

	})

	r.POST("money/createCharge", func(c *gin.Context) {
		var chargeInstance types.Charge
		err := c.BindJSON(&chargeInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		chargeInstance.Amount = chargeInstance.Amount * 100
		// `source` is obtained with Stripe.js; see https://stripe.com/docs/payments/accept-a-payment-charges#web-create-token
		params := &stripe.ChargeParams{
			Amount:   stripe.Int64(chargeInstance.Amount),
			Currency: stripe.String(string(chargeInstance.Currency)),
			Customer: stripe.String(chargeInstance.Customer),
		}
		charge, _ := charge.New(params)

		c.JSON(200, gin.H{
			"charge": charge,
		})

	})

	r.POST("money/createTransfer", func(c *gin.Context) {
		var transferInstance types.Transfer
		err := c.BindJSON(&transferInstance)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}

		transferInstance.Amount = transferInstance.Amount * 100

		params := &stripe.TransferParams{
			Amount:      stripe.Int64(transferInstance.Amount),
			Currency:    stripe.String(string(transferInstance.Currency)),
			Destination: stripe.String(transferInstance.Dest),
		}
		transfer, _ := transfer.New(params)

		c.JSON(200, gin.H{
			"transfer": transfer,
		})

	})
}

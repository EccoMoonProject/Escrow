package types

import "github.com/stripe/stripe-go/v74"

type Product struct {
	Name  string `bson:"name"`
	Price int64  `bson:"price"`
}

type Price struct {
	Product  string `bson:"product"`
	Amount   int64  `bson:"amount"`
	Currency string `bson:"currency"`
}

type Customer struct {
	Email   string `bson:"email"`
	OwnerID string `bson:"ownerID"`
}

type PaymentMethod struct {
	CustomerID string `bson:"customerID"`
	Type       string `bson:"type"`
	Card       Card   `bson:"card"`
}

type Card struct {
	Number   string `bson:"number"`
	ExpMonth int64  `bson:"expMonth"`
	ExpYear  int64  `bson:"expYear"`
	CVC      string `bson:"cvc"`
}

type Charge struct {
	Amount   int64  `bson:"amount"`
	Currency string `bson:"currency"`
	Customer string `bson:"customer"`
}

type Transfer struct {
	Amount   int64  `bson:"amount"`
	Currency string `bson:"currency"`
	Dest     string `bson:"dest"`
}

type CustomerPortal struct {
	Customer string `bson:"customer"`
}

type StripeAccount struct {
	Country string `bson:"country"`
	Email   string `bson:"email"`
}

type Payouts struct {
	Amount      int64  `bson:"amount"`
	Currency    string `bson:"currency"`
	Destination string `bson:"destination"`
}

type BankAccount struct {
	Customer          string `bson:"customer"`
	AccountHolderName string `bson:"accountHolderName"`
	Country           string `bson:"country"`
	Currency          string `bson:"currency"`
	RoutingNumber     string `bson:"routingNumber"`
	AccountNumber     string `bson:"accountNumber"`
}

type CreateCardRequest struct {
	Customer string `bson:"customer"`
	Name     string `bson:"name"`
	Currency string `bson:"currency"`
	Number   string `bson:"number"`
	ExpMonth string `bson:"expMonth"`
	ExpYear  string `bson:"expYear"`
	CVC      string `bson:"cvc"`
}

type CardData struct {
	Card     *stripe.Card `bson:"card"`
	Customer string       `bson:"customer"`
}

type PayInstanceRequest struct {
	InstanceID string `bson:"instanceID"`
	OwnerID    string `bson:"ownerID"`
	OwnerSHI   string `bson:"ownerSHI"`
	Amount     uint64 `bson:"amount"`
}

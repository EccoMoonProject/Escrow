package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name     string `bson:"name"`
	OwnerID  string `bson:"ownerID"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type EscrowInstance struct {
	InstanceID      string `bson:"instanceID"`
	OwnerID         string `bson:"ownerID"`
	OwnerEmail      string `bson:"ownerEmail"`
	Amount          uint64 `bson:"amount"`
	Status          bool   `bson:"status"`
	OwnerSHI        string `bson:"ownerSHI"`
	SecureDestroyer bool   `bson:"secureDestroyer"`
}

type EscrowRequest struct {
	OwnerID  string `bson:"ownerID"`
	Email    string `bson:"email"`
	Amount   uint64 `bson:"amount"`
	Currency string `bson:"currency"`
}

type VotingPool struct {
	InstanceID string `bson:"instanceID"`
	BuyerID    string `bson:"buyerID"`
	SellerID   string `bson:"sellerID"`
	BuyerVote  uint16 `bson:"buyerVote"`
	SellerVote uint16 `bson:"sellerVote"`
	Consensus  bool   `bson:"consensus"`
}

type VotingRequest struct {
	InstanceID string `bson:"instanceID"`
	BuyerID    string `bson:"ownerID"`
	SellerID   string `bson:"sellerID"`
}

type PaymentRequest struct {
	InstanceID string `bson:"instanceID"`
	OwnerID    string `bson:"ownerID"`
	OwnerSHI   string `bson:"ownerSHI"`
	Amount     uint64 `bson:"amount"`
}

type Wallet struct {
	OwnerID string `bson:"ownerID"`
	Balance uint64 `bson:"balance"`
}

type DepositRequest struct {
	OwnerID string `bson:"ownerID"`
	Amount  uint64 `bson:"amount"`
}

type WithdrawRequest struct {
	OwnerID string `bson:"ownerID"`
	Amount  uint64 `bson:"amount"`
}

type Dispute struct {
	OwnerID  string `bson:"ownerID"`
	OwnerSHI string `bson:"ownerSHI"`
	Category string `bson:"category"`
	Amount   uint64 `bson:"amount"`
}

// Define a struct to represent a chat
type Chat struct {
	ID        string `bson:"_id"`
	Sender    string `bson:"sender"`
	OwnerSHI  string `bson:"ownerSHI"`
	Recipient string `bson:"recipient"`
	Message   string `bson:"message"`
}

type ChatRoom struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Members []string           `bson:"members"`
}

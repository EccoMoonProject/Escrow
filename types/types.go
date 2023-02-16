package types

type User struct {
	Name        string `bson:"name"`
	Email       string `bson:"email"`
	Password    string `bson:"password"`
	DateOfBirth string `bson:"dateOfBirth"`
}

type EscrowInstance struct {
	InstanceID      string `bson:"instanceID"`
	OwnerID         string `bson:"ownerID"`
	OwnerName       string `bson:"ownerName"`
	OwnerEmail      string `bson:"ownerEmail"`
	OwnerPhone      string `bson:"ownerPhone"`
	Amount          uint64 `bson:"amount"`
	Status          bool   `bson:"status"`
	OwnerSHI        string `bson:"ownerSHI"`
	SecureDestroyer bool   `bson:"secureDestroyer"`
}

type VotingPool struct {
	InstanceID string `bson:"instanceID"`
	BuyerID    string `bson:"buyerID"`
	SellerID   string `bson:"sellerID"`
	BuyerVote  uint16 `bson:"buyerVote"`
	SellerVote uint16 `bson:"sellerVote"`
	Consensus  bool   `bson:"consensus"`
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

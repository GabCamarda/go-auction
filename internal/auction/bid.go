package auction

import (
	"time"
)

// Bid represents a user bid
type Bid struct {
	userID     int
	itemID     string
	amount     int
	occurredAt time.Time
}

func NewBid(userID, amount int, itemID string, occurredAt time.Time) (Bid, error) {
	// depending on the source, a userID with value 0 could be valid
	if userID < 0 {
		return Bid{}, ErrInvalidBidUserID
	}
	if amount <= 0 {
		return Bid{}, ErrInvalidBidAmount
	}
	if itemID == "" {
		return Bid{}, ErrInvalidBidItemID
	}
	if occurredAt.IsZero() {
		return Bid{}, ErrInvalidBidOccurredAt
	}

	return Bid{
		userID:     userID,
		itemID:     itemID,
		amount:     amount,
		occurredAt: occurredAt,
	}, nil
}

func (b Bid) UserID() int {
	return b.userID
}

func (b Bid) ItemID() string {
	return b.itemID
}

func (b Bid) Amount() int {
	return b.amount
}

func (b Bid) OccurredAt() time.Time {
	return b.occurredAt
}

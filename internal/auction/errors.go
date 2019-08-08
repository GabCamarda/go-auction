package auction

import "errors"

// Common error definitions
var (
	// Bid
	ErrInvalidBidUserID     = errors.New("bid: invalid user id")
	ErrInvalidBidItemID     = errors.New("bid: invalid item id")
	ErrInvalidBidAmount     = errors.New("bid: invalid amount, must be greater than zero")
	ErrInvalidBidOccurredAt = errors.New("bid: invalid time")

	// Item
	ErrInvalidItemSellerID     = errors.New("item: invalid seller id")
	ErrInvalidItemID           = errors.New("item: invalid item id")
	ErrInvalidItemReservePrice = errors.New("item: invalid reserved price, must be greater than zero")

	// Auction
	ErrInvalidAuctionStartTime   = errors.New("auction: invalid start time")
	ErrInvalidAuctionCloseTime   = errors.New("auction: invalid close time")
	ErrClosedAuctionBid          = errors.New("auction: bidding is no longer accepted, auction is closed")
	ErrAuctionBidNotHigherAmount = errors.New("auction: bid amount must be greater than the highest bid")

	// Service
	ErrServiceCannotStore = errors.New("service: the auction could not be stored, invalid id")
	ErrServiceNotFound    = errors.New("service: auction not found")
)

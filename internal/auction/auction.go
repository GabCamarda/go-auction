package auction

import (
	"sort"
	"time"
)

type Status int

const (
	sold Status = iota
	unsold
)

func (s Status) ToString() string {
	switch s {
	case sold:
		return "SOLD"
	default:
		return "UNSOLD"
	}
}

// Winner represents the winning bid
// It contains the highest bid and the final amount to pay
// The amount to pay is selected as per the following:
// - if there's only 1 bid, the amount equals to item reserve price
// - if more than 1 bid, the amount equals the second highest bid
type Winner struct {
	bid           Bid
	itemPriceSold int
}

// Bid returns the winning bid
func (w Winner) Bid() Bid {
	return w.bid
}

// ItemPriceSold returns the price the item was sold as
func (w Winner) ItemPriceSold() int {
	return w.itemPriceSold
}

// Auction represent the auction entity
type Auction struct {
	sellingItem Item
	bids        []Bid
	highestBid  Bid
	lowestBid   Bid
	winner      Winner
	startTime   time.Time
	closeTime   time.Time
	status      Status
	closed      bool
}

// NewAuction returns a pointer to Auction, error in case of validation errors
func NewAuction(sellingItem Item, startTime, closeTime time.Time) (*Auction, error) {
	if startTime.IsZero() {
		return nil, ErrInvalidAuctionStartTime
	}
	if closeTime.IsZero() {
		return nil, ErrInvalidAuctionCloseTime
	}

	return &Auction{
		sellingItem: sellingItem,
		bids:        []Bid{}, // slices can be nil
		startTime:   startTime,
		closeTime:   closeTime,
		status:      unsold,
		closed:      false,
	}, nil
}

// Bid adds a bid to the auction
// if the bid occurred after the auction has closed
// it will return an error
func (a *Auction) Bid(bid Bid) error {
	if a.status == sold || a.closed {
		return ErrClosedAuctionBid
	}
	if bid.OccurredAt().After(a.closeTime) {
		if !a.closed {
			a.Close(bid.OccurredAt())
		}
		return ErrClosedAuctionBid
	}

	highestBid, ok := a.HighestBid()
	if ok {
		if bid.Amount() <= highestBid.Amount() {
			return ErrAuctionBidNotHigherAmount
		}
	}

	a.bids = append(a.bids, bid)
	return nil
}

func (a *Auction) sortBidsByAmount() {
	sort.Slice(a.bids, func(i, j int) bool {
		return a.bids[i].Amount() < a.bids[j].Amount()
	})
}

// HighestBid returns the highest bid or false if no bids have been submitted
func (a *Auction) HighestBid() (Bid, bool) {
	if a.TotalBidCount() == 0 {
		return Bid{}, false
	}
	if a.highestBid != (Bid{}) {
		return a.highestBid, true
	}

	a.sortBidsByAmount()
	a.highestBid = a.bids[len(a.bids)-1]

	return a.highestBid, true
}

// LowestBid returns the lowest bid or false if no bids have been submitted
func (a *Auction) LowestBid() (Bid, bool) {
	if a.TotalBidCount() == 0 {
		return Bid{}, false
	}
	if a.lowestBid != (Bid{}) {
		return a.lowestBid, true
	}

	a.sortBidsByAmount()
	a.lowestBid = a.bids[0]

	return a.lowestBid, true
}

// TotalBidCount returns the total number of bids
func (a *Auction) TotalBidCount() int {
	return len(a.bids)
}

func (a *Auction) canBeClosed(at time.Time) bool {
	return a.CloseTime().Before(at) || a.CloseTime().Equal(at)
}

// Close closes the auction if the time argument specified is equal
// or before the auction close time
// If it can be closed, a winner is chosen and the right amount to pay is calculated
// Returns false if the auction cannot be closed, true if it has closed
func (a *Auction) Close(at time.Time) bool {
	if !a.canBeClosed(at) || a.closed {
		return false
	}

	switch a.TotalBidCount() {
	case 0: // no bids, item is unsold
		a.status = unsold
		a.closed = true
		return true
	case 1: // only 1 bid, item is sold at the item reserve price
		a.sortBidsByAmount()
		a.highestBid = a.bids[1]
		a.lowestBid = a.bids[0]
		return a.createWinningBid(a.sellingItem.ReservePrice())
	default: // item is sold at the second highest bid amount
		a.sortBidsByAmount()
		a.highestBid = a.bids[len(a.bids)-1]
		a.lowestBid = a.bids[0]
		return a.createWinningBid(a.bids[len(a.bids)-2].Amount())
	}
}

func (a *Auction) createWinningBid(itemPriceSold int) bool {
	if a.highestBid.Amount() < a.sellingItem.ReservePrice() {
		a.status = unsold
		a.closed = true
		return true
	}
	a.winner = Winner{bid: a.highestBid, itemPriceSold: itemPriceSold}
	a.status = sold
	a.closed = true
	return true
}

// ID returns the auction id, which equals the selling item unique id
func (a *Auction) ID() string {
	return a.sellingItem.ID()
}

// Winner returns the auction winner
func (a *Auction) Winner() Winner {
	return a.winner
}

// StartTime returns the time when the auction started
func (a *Auction) StartTime() time.Time {
	return a.startTime
}

// CloseTime returns the time when the auction ended
func (a *Auction) CloseTime() time.Time {
	return a.closeTime
}

func (a *Auction) Status() string {
	return a.status.ToString()
}

package auction_test

import (
	"testing"
	"time"

	"github.com/GabCamarda/go-auction/internal/auction"
)

func TestAuctionSuccess(t *testing.T) {
	itemID := "item id"
	sellerID := 1
	userID1 := 3
	userID2 := 6
	item, _ := auction.NewItem(itemID, sellerID, 2500)

	start := time.Now().UTC()
	closeIn, _ := time.ParseDuration("15s")
	close := start.Add(closeIn)
	a, _ := auction.NewAuction(item, start, close)

	bid1After, _ := time.ParseDuration("5s")
	bidAt := start.Add(bid1After)
	bid1, _ := auction.NewBid(userID1, 1300, itemID, bidAt)

	bid2After, _ := time.ParseDuration("10s")
	bid2At := start.Add(bid2After)
	bid2, _ := auction.NewBid(userID2, 1800, itemID, bid2At)

	bid3After, _ := time.ParseDuration("15s")
	bid3At := start.Add(bid3After)
	bid3, _ := auction.NewBid(userID1, 2800, itemID, bid3At)

	a.Bid(bid1)
	a.Bid(bid2)
	a.Bid(bid3)
	if !a.Close(close) {
		t.Error("auction should close without error")
	}
	if a.TotalBidCount() != 3 {
		t.Error("auction should have 3 bids")
	}
	highestBid, _ := a.HighestBid()
	if highestBid != bid3 {
		t.Errorf("bid3 should be the highest, found: %v", highestBid)
	}
	lowestBid, _ := a.LowestBid()
	if lowestBid != bid1 {
		t.Errorf("bid1 should be the lowest, found: %v", lowestBid)
	}
	winner := a.Winner()
	if winner.ItemPriceSold() != bid2.Amount() {
		t.Error("winner should pay the second highest bid amount")
	}

	if a.ID() != itemID {
		t.Errorf("auction id and item id should be equal, auctionID: %s, itemID: %s", a.ID(), itemID)
	}
	if a.Status() != "SOLD" {
		t.Errorf("auction status should be SOLD, found: %s", a.Status())
	}
}

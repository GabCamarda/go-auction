package auction

// Item represents the listing item for the auction
// The reserved price is the seller initial selling price
type Item struct {
	id           string
	sellerID     int
	reservePrice int
}

func NewItem(id string, sellerID, reservePrice int) (Item, error) {
	if id == "" {
		return Item{}, ErrInvalidItemID
	}
	if sellerID < 1 {
		return Item{}, ErrInvalidItemSellerID
	}
	if reservePrice <= 0 {
		return Item{}, ErrInvalidItemReservePrice
	}

	return Item{
		id:           id,
		sellerID:     sellerID,
		reservePrice: reservePrice,
	}, nil
}

func (i Item) ID() string {
	return i.id
}

func (i Item) SellerID() int {
	return i.sellerID
}

func (i Item) ReservePrice() int {
	return i.reservePrice
}

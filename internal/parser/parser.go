package parser

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/GabCamarda/go-auction/internal/auction"
)

const (
	auctionSell = "SELL"
	auctionBid  = "BID"
)

type AuctionService interface {
	Store(a *auction.Auction) error
	GetByItemID(itemID string) (*auction.Auction, error)
	CompletedAuctions(at time.Time) []*auction.Auction
}

type Writer interface {
	Write(a *auction.Auction)
}

// Parser is used to parse auction data from a reader
type Parser struct {
	auctionService AuctionService
	writer         Writer
}

// New returns a pointer to Parser
func New(auctionService AuctionService, writer Writer) *Parser {
	return &Parser{
		auctionService: auctionService,
		writer:         writer,
	}
}

// Parse parses auction data from a reader
// The client is responsible for handling file connections
func (p *Parser) Parse(reader *bufio.Reader) error {
	if reader == nil {
		return errors.New("reader cannot be nil")
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, auctionSell) {
			err := p.parseAuctionSelling(s)
			if err != nil {
				return err
			}
		} else if strings.Contains(s, auctionBid) {
			err := p.parseAuctionBid(s)
			if err != nil {
				return err
			}
		} else if s != "" {
			p.heartBeat(s)
		}
	}

	return scanner.Err()
}

// parses selling input lines
func (p *Parser) parseAuctionSelling(fileLine string) error {
	chunks := strings.Split(fileLine, "|")
	startTime, err := p.stringToTime(chunks[0])
	if err != nil {
		return err
	}
	closeTime, err := p.stringToTime(chunks[5])
	if err != nil {
		return err
	}
	reservePrice, err := p.stringToMoney(chunks[4])
	if err != nil {
		return err
	}
	userID, err := strconv.Atoi(chunks[1])
	if err != nil {
		return err
	}
	itemID := chunks[3]

	item, err := auction.NewItem(itemID, userID, reservePrice)
	if err != nil {
		return err
	}

	a, err := auction.NewAuction(item, startTime, closeTime)
	if err != nil {
		return err
	}

	p.storeAuction(a)
	return nil
}

// parses bid input lines
func (p *Parser) parseAuctionBid(fileLine string) error {
	chunks := strings.Split(fileLine, "|")
	occurredTime, err := p.stringToTime(chunks[0])
	if err != nil {
		return err
	}
	bidAmount, err := p.stringToMoney(chunks[4])
	if err != nil {
		return err
	}
	userID, err := strconv.Atoi(chunks[1])
	if err != nil {
		return err
	}
	itemID := chunks[3]

	bid, err := auction.NewBid(userID, bidAmount, itemID, occurredTime)
	if err != nil {
		return err
	}

	p.storeBid(bid)
	return nil
}

// converts a string to a timestamp
func (p *Parser) stringToTime(val string) (time.Time, error) {
	timeInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(timeInt, 0).UTC(), nil
}

// converts a string to currency
// assumption: the format is always with 2 digit precision
// ints are used internally to handle money
func (p *Parser) stringToMoney(val string) (int, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	// return amount in cents
	return int(f * 100), nil
}

func (p *Parser) storeAuction(auction *auction.Auction) error {
	if err := p.auctionService.Store(auction); err != nil {
		return err
	}

	return nil
}

func (p *Parser) storeBid(bid auction.Bid) error {
	auction, err := p.auctionService.GetByItemID(bid.ItemID())
	if err != nil {
		return err
	}

	// for this exercise, ignore bids submitted after the auction was closed
	// in a real word scenario, this should probably be logged, returned to the client as an error
	// or stored in a database for auditing purposes
	auction.Bid(bid)

	return p.auctionService.Store(auction)
}

// checks if any auctions can be closed
func (p *Parser) heartBeat(fileLine string) error {
	chunks := strings.Split(fileLine, "|")
	at, err := p.stringToTime(chunks[0])
	if err != nil {
		return err
	}

	auctions := p.auctionService.CompletedAuctions(at)
	for _, auction := range auctions {
		p.writer.Write(auction)
	}

	return nil
}

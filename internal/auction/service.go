package auction

import (
	"log"
	"sort"
	"time"

	"github.com/GabCamarda/go-auction/internal/utils/logutil"
)

// Service is the auction service controller
// For brevity a map is used to store auctions
type Service struct {
	// maps are not concurrent safe, but for the purpose of this exercise
	// a mutex is not needed
	data   map[string]*Auction
	logger *log.Logger
}

func NewService(logger *log.Logger) *Service {
	if logger == nil {
		logger = logutil.NewDefaultLogger()
	}

	return &Service{
		logger: logger,
		data:   make(map[string]*Auction),
	}
}

func (s *Service) Store(a *Auction) error {
	if a.ID() == "" {
		return ErrServiceCannotStore
	}

	s.data[a.ID()] = a
	return nil
}

func (s *Service) GetByItemID(itemID string) (*Auction, error) {
	auction, ok := s.data[itemID]
	if !ok {
		return nil, ErrServiceNotFound
	}

	return auction, nil
}

func (s *Service) GetAll() []*Auction {
	var auctions []*Auction
	for _, auction := range s.data {
		auctions = append(auctions, auction)
	}

	return auctions
}

func (s *Service) CompletedAuctions(at time.Time) []*Auction {
	auctions := s.GetAll()
	sort.Slice(auctions, func(i, j int) bool {
		return auctions[i].StartTime().Before(auctions[j].StartTime())
	})

	var completedAuctions []*Auction
	for _, auction := range auctions {
		if auction.Close(at) {
			completedAuctions = append(completedAuctions, auction)
		}
	}

	return completedAuctions
}

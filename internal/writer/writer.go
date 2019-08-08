package writer

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/GabCamarda/go-auction/internal/auction"
)

type Writer struct {
	writer io.Writer
}

func New(writer io.Writer) *Writer {
	if writer == nil {
		writer = os.Stdout
	}

	return &Writer{
		writer: writer,
	}
}

func (w *Writer) Write(a *auction.Auction) {
	winner := a.Winner()
	userID := ""
	if winner.Bid().UserID() != 0 {
		userID = strconv.Itoa(winner.Bid().UserID())
	}

	highestBidAmount := "0.00"
	lowestBidAmount := "0.00"
	highestBid, ok := a.HighestBid()
	if ok {
		highestBidAmount = fmt.Sprintf("%.2f", float64(highestBid.Amount())/100)
	}
	lowestBid, ok := a.LowestBid()
	if ok {
		lowestBidAmount = fmt.Sprintf("%.2f", float64(lowestBid.Amount())/100)
	}

	// a bit verbose, but it's the only way to display this in order without weird hacks
	// a better way is to use a struct, but it also prints {}
	w.formatWrite("%v|", a.CloseTime().Unix())
	w.formatWrite("%v|", a.ID())
	w.formatWrite("%v|", userID)
	w.formatWrite("%v|", a.Status())
	w.formatWrite("%v|", fmt.Sprintf("%.2f", float64(a.Winner().ItemPriceSold())/100))
	w.formatWrite("%v|", a.TotalBidCount())
	w.formatWrite("%v|", highestBidAmount)
	w.formatWrite("%v", lowestBidAmount)
	w.writer.Write([]byte("\n"))
}

func (w *Writer) formatWrite(format string, val interface{}) {
	w.writer.Write([]byte(fmt.Sprintf(format, val)))
}

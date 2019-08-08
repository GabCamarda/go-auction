package parser_test

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/GabCamarda/go-auction/internal/auction"
	"github.com/GabCamarda/go-auction/internal/parser"
	"github.com/GabCamarda/go-auction/internal/utils/logutil"
	"github.com/GabCamarda/go-auction/internal/writer"
)

func TestParser(t *testing.T) {
	file, err := os.Open("./test/testinput.txt")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	logger := logutil.NewDefaultLogger()
	auctionService := auction.NewService(logger)
	var b bytes.Buffer
	buf := bufio.NewWriter(&b)
	err = parser.New(auctionService, writer.New(buf)).Parse(bufio.NewReader(file))
	if err != nil {
		t.Fail()
	}

	buf.Flush()
	out := b.String()
	result1 := "20|toaster_1|8|SOLD|12.50|3|20.00|7.50"
	result2 := "20|tv_1||UNSOLD|0.00|2|200.00|150.00"
	if !strings.Contains(out, result1) {
		t.Error("parsing result is not as expected")
	}
	if !strings.Contains(out, result2) {
		t.Error("parsing result is not as expected")
	}
}

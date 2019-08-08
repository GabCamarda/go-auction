package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/GabCamarda/go-auction/internal/auction"
	"github.com/GabCamarda/go-auction/internal/parser"
	"github.com/GabCamarda/go-auction/internal/utils/logutil"
	"github.com/GabCamarda/go-auction/internal/writer"
)

var filepath = flag.String("filepath", "", "the filepath where the data is stored")

func main() {
	flag.Parse()

	file, err := os.Open(*filepath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	logger := logutil.NewDefaultLogger()
	auctionService := auction.NewService(logger)
	err = parser.New(auctionService, writer.New(os.Stdout)).Parse(bufio.NewReader(file))
	if err != nil {
		fmt.Println(err)
	}
}

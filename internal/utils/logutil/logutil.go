package logutil

import (
	"log"
	"os"
)

func NewDefaultLogger() *log.Logger {
	return log.New(os.Stdout, "tm-auction-service: ", log.LstdFlags|log.LUTC)
}

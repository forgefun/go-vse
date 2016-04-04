package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/forgefun/go-vse/vse"
	"github.com/joho/godotenv"
)

func main() {
	log.SetLevel(log.DebugLevel)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := vse.NewClient(vse.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	portfolio := client.Portfolio("sim101")
	portfolio.GetHoldings()

	order := &vse.Order{
		Fuid:   "STOCK-XNAS-AAPL",
		Shares: "1",
		Type:   "Buy",
		Term:   "Cancelled",
	}

	portfolio.SubmitOrder(*order)

	portfolio.ListOrders()
}

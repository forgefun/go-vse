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

	err = client.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	portfolio := client.Portfolio("sim101")
	if _, err := portfolio.GetHoldings(); err != nil {
		log.Fatal(err)
	}

	order := &vse.Order{
		Fuid:   "STOCK-XASQ-AA",
		Shares: "1",
		Type:   "Buy",
		Term:   "Cancelled",
	}

	if err := portfolio.SubmitOrder(*order); err != nil {
		log.Fatal(err)
	}

	// if err := portfolio.CancelOrder("80916009"); err != nil {
	// 	log.Fatal(err)
	// }

	if _, err := portfolio.ListOrders(); err != nil {
		log.Fatal(err)
	}
}

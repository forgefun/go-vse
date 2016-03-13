package main

import (
  "log"
  // "os"

  "github.com/cleung2010/go-vse/vse"
  "github.com/joho/godotenv"
)

func main()  {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  client, err := vse.NewClient(vse.DefaultConfig())
  if err != nil {
    log.Fatal(err)
  }

  portfolio := client.Portfolio("sim101")
  portfolio.GetHoldings()
}

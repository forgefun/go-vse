package main

import (
  "log"
  "os"

  "github.com/cleung2010/go-vse/vse"
  "github.com/joho/godotenv"
)

func main()  {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  username := os.Getenv("USERNAME")
  password := os.Getenv("PASSWORD")

  vse.Authenticate(username, password)
}

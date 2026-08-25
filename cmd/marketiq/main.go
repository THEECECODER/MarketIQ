package main

import (
	"log"

	"marketiq/internal/marketdata"
)

func main() {
	log.Println("Starting MarketIQ...")

	client := marketdata.NewClient()

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	// Keep the application alive.
	select {}
}

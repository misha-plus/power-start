package main

import (
	"log"
	"net/http"
	"time"
)

func sendBeacon() {
	// TODO extract url
	// TODO timeout
	resp, err := http.DefaultClient.Post(
		"http://localhost:3000/agent/thename/heartbeat", "", nil)
	if err != nil {
		log.Printf("Heartbeat error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Heartbeat error: status = %d", resp.StatusCode)
	}
}

func main() {
	// TODO extract interval
	ticker := time.NewTicker(5 * time.Second)
	log.Println("Started")
	for range ticker.C {
		sendBeacon()
	}
}

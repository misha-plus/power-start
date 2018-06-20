package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func sendBeacon() {
	// TODO extract url
	// TODO timeout
	// TODO auth
	resp, err := http.DefaultClient.Post(
		"http://localhost:3000/agent/thename/heartbeat", "", nil)
	if err != nil {
		log.Printf("Heartbeat error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Heartbeat error: status = %d", resp.StatusCode)
		return
	}

	data := struct {
		ShouldShutdown bool `json:"shouldShutdown"`
	}{false}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	if data.ShouldShutdown {
		shutdown()
	}
}

func shutdown() {
	println("=====================")
	println("Shutting down")
	println("=====================")
}

func main() {
	// TODO extract interval
	ticker := time.NewTicker(5 * time.Second)
	log.Println("Started")
	for range ticker.C {
		sendBeacon()
	}
}

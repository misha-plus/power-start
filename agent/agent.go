package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

var config = struct {
	ServerBaseURL   string
	MachineName     string
	ShutdownCommand string
	Heartbeat       struct {
		IntervalSeconds int
		TimeoutSeconds  int
	}
}{}

var client http.Client

func sendBeacon() {
	// TODO auth
	url := fmt.Sprintf(
		"%s/api/agent/%s/heartbeat", config.ServerBaseURL, config.MachineName)
	resp, err := client.Post(url, "", nil)
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

	if !data.ShouldShutdown {
		return
	}

	err = shutdown()
	if err != nil {
		fmt.Printf("Can't shutdown: %v", err)
	}
}

func shutdown() error {
	fmt.Printf("Shutting down using '%s'\n", config.ShutdownCommand)
	cmd := exec.Command("bash", "-c", config.ShutdownCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	byline := strings.Split(string(output), "\n")
	for _, line := range byline {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fmt.Println(">", line)
	}

	return nil
}

func main() {
	err := ini.MapTo(&config, "agent-config.ini")
	if err != nil {
		log.Fatalf("Can't parse config: %v", err)
	}

	client = http.Client{
		Timeout: time.Duration(config.Heartbeat.TimeoutSeconds) * time.Second,
	}

	ticker := time.NewTicker(
		time.Duration(config.Heartbeat.IntervalSeconds) * time.Second)
	log.Println("Started")
	sendBeacon()
	for range ticker.C {
		sendBeacon()
	}
}

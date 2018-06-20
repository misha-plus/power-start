package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/go-ini/ini"
)

var serverBaseURL = "http://localhost:3000"
var machineName = "thename"
var interval = 5 * time.Second

var config = struct {
	ServerBaseURL   string
	MachineName     string
	IntervalSeconds int
}{}

func sendBeacon() {
	// TODO timeout
	// TODO auth
	url := fmt.Sprintf(
		"%s/agent/%s/heartbeat", config.ServerBaseURL, config.MachineName)
	resp, err := http.DefaultClient.Post(url, "", nil)
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
	fmt.Println("Shutting down")
	cmd := exec.Command("echo", "Hello")

	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	defer stderr.Close()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	return cmd.Run()
}

func main() {
	err := ini.MapTo(&config, "agent-config.ini")
	if err != nil {
		log.Fatalf("Can't parse config: %v", err)
	}

	ticker := time.NewTicker(time.Duration(config.IntervalSeconds) * time.Second)
	log.Println("Started")
	for range ticker.C {
		sendBeacon()
	}
}

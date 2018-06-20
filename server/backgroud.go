package main

import (
	"encoding/json"
	"log"
	"time"

	bolt "github.com/coreos/bbolt"
)

func (h *appHandle) backgroundJob() {
	err := h.db.View(func(tx *bolt.Tx) error {
		machines := tx.Bucket(machineBucket)
		c := machines.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			machine := machineRecord{}
			err := json.Unmarshal(v, &machine)
			if err != nil {
				log.Printf("Background job/take machine: %v", err)
				continue
			}

			isRunning := time.Now().Sub(machine.LastHeartbeat) < 5*time.Minute
			if isRunning || machine.Requests == 0 {
				continue
			}

			err = startMachine(machine.MAC)
			if err != nil {
				log.Printf("Background job/start machine: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Background job: %v", err)
	}
}

func (h *appHandle) runBackgoundJobs() {
	d := 30 * time.Second
	timer := time.NewTimer(d)
	for {
		<-timer.C
		h.backgroundJob()
		timer.Reset(d)
	}
}

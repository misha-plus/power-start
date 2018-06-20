package main

import (
	"encoding/json"
	"time"

	bolt "github.com/coreos/bbolt"
)

type machineRecord struct {
	Name          string `json:"name"`
	MAC           string `json:"mac"`
	Requests      int    // TODO: hide it from output
	LastHeartbeat time.Time
	LastRequest   time.Time
}

var machineBucket = []byte("Machines")

func getMachine(tx *bolt.Tx, name string) (*machineRecord, error) {
	machines := tx.Bucket(machineBucket)
	machine := &machineRecord{}
	bytes := machines.Get([]byte(name))
	if bytes == nil {
		return nil, nil
	}
	err := json.Unmarshal(bytes, machine)
	if err != nil {
		return nil, err
	}
	return machine, nil
}

func putMachine(tx *bolt.Tx, machine *machineRecord) error {
	bytes, err := json.Marshal(machine)
	if err != nil {
		return err
	}
	machines := tx.Bucket(machineBucket)
	return machines.Put([]byte(machine.Name), bytes)
}

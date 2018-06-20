package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	bolt "github.com/coreos/bbolt"
	"github.com/go-ini/ini"
)

var config = struct {
	MachineInactivityTimeoutSeconds int
	BindAddress                     string
	DBPath                          string
}{}

type appHandle struct {
	db *bolt.DB
}

func newApp() (*appHandle, error) {
	db, err := bolt.Open(config.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(machineBucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	log.Println("DB is prepared")

	return &appHandle{db}, nil
}

func (app *appHandle) stop() {
	app.db.Close()
}

func main() {
	err := ini.MapTo(&config, "server-config.ini")
	if err != nil {
		log.Fatalf("Can't parse config: %v", err)
	}

	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.stop()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		app.stop()
		os.Exit(0)
	}()

	go app.runBackgoundJobs()
	app.runRest()
}

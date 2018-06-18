package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	bolt "github.com/coreos/bbolt"
	"github.com/go-chi/chi"
	"github.com/linde12/gowol"
)

type record struct {
	Name      string `json:"name"`
	MAC       string `json:"mac"`
	IsRunning bool
}

var machineBucket = []byte("Machines")

type httpError struct {
	code    int
	message string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("%d - %s", e.code, e.message)
}

func newHTTPError(code int, message string) *httpError {
	return &httpError{code, message}
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	defer log.Print("DB is closed")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		db.Close()
		log.Print("DB is closed")
		os.Exit(0)
	}()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(machineBucket)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB is prepared")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		theRecord := record{}
		err := json.NewDecoder(r.Body).Decode(&theRecord)
		if err != nil {
			log.Printf("/add: %v", err)
			w.WriteHeader(400)
			w.Write([]byte("Invalid JSON"))
			return
		}

		err = db.Update(func(tx *bolt.Tx) error {
			machines := tx.Bucket(machineBucket)
			data, err := json.Marshal(theRecord)
			if err != nil {
				return err
			}
			return machines.Put([]byte(theRecord.Name), []byte(data))
		})
		if err != nil {
			log.Printf("/add: %v", err)
			w.WriteHeader(500)
			w.Write([]byte("Internal error"))
			return
		}

		log.Printf("Added machine: %v", theRecord)
		// TODO: add validation
		w.Write([]byte("Ok"))
	})

	r.Post("/remove/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		err := db.Update(func(tx *bolt.Tx) error {
			machines := tx.Bucket(machineBucket)
			return machines.Delete([]byte(name))
		})
		if err != nil {
			log.Printf("/remove/%s: %v", name, err)
			w.WriteHeader(500)
			w.Write([]byte("Internal error"))
			return
		}
		w.Write([]byte("Deleted record for " + name))
	})

	r.Post("/start/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		err := db.View(func(tx *bolt.Tx) error {
			machine := tx.Bucket(machineBucket)
			theRecord := record{}
			json.Unmarshal(machine.Get([]byte(name)), &theRecord)

			packet, err := gowol.NewMagicPacket(theRecord.MAC)
			if err != nil {
				return err
			}
			// TODO: add selecting port and IP
			log.Printf("Starting machine: %s", name)
			return packet.Send("255.255.255.255")
		})

		if err != nil {
			log.Printf("/start/%s: %v", name, err)
			w.WriteHeader(500)
			w.Write([]byte("Internal error"))
			return
		}

		w.Write([]byte("Starting"))
	})

	r.Get("/status/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		var result []record
		err := db.View(func(tx *bolt.Tx) error {
			machines := tx.Bucket(machineBucket)
			c := machines.Cursor()
			var err error

			for k, v := c.First(); k != nil && err == nil; k, v = c.Next() {
				theRecord := record{}
				err = json.Unmarshal(v, &theRecord)
				result = append(result, theRecord)
			}
			return err
		})

		if err != nil {
			log.Printf("/list: %v", err)
			w.WriteHeader(500)
			w.Write([]byte("Internal error"))
			return
		}

		json.NewEncoder(w).Encode(result)
	})

	// TODO add auth
	r.Post("/agent/{name}/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		err := db.Batch(func(tx *bolt.Tx) error {
			machine := tx.Bucket(machineBucket)
			theRecord := record{}
			bytes := machine.Get([]byte(name))
			if bytes == nil {
				return newHTTPError(404, "Machine not found")
			}
			json.Unmarshal(bytes, &theRecord)

			if theRecord.IsRunning {
				return nil
			}

			theRecord.IsRunning = true
			bytes, err := json.Marshal(theRecord)
			if err != nil {
				return err
			}

			fmt.Println(theRecord)

			return machine.Put([]byte(theRecord.Name), bytes)
		})

		if err, ok := err.(*httpError); ok {
			log.Printf("/agent/%s/heartbeat: %v", name, err)
			w.WriteHeader(err.code)
			w.Write([]byte(err.message))
			return
		}

		if err != nil {
			log.Printf("/agent/%s/heartbeat: %v", name, err)
			w.WriteHeader(500)
			w.Write([]byte("Internal error"))
			return
		}

		w.Write([]byte("Ok"))
	})
	http.ListenAndServe(":3000", r)
	log.Println("main")
}

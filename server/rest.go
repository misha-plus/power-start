package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/go-chi/chi"
)

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

func respondError(err error, w http.ResponseWriter, r *http.Request) {
	if err, ok := err.(*httpError); ok {
		log.Printf("%s: %v", r.RequestURI, err)
		w.WriteHeader(err.code)
		w.Write([]byte(err.message))
		return
	}

	if err != nil {
		log.Printf("%s: %v", r.RequestURI, err)
		w.WriteHeader(500)
		w.Write([]byte("Internal error"))
	}
}

func (app *appHandle) runRest() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		machine := machineRecord{}
		err := json.NewDecoder(r.Body).Decode(&machine)
		if err != nil {
			respondError(newHTTPError(400, "Invalid JSON"), w, r)
			return
		}
		if _, err := net.ParseMAC(machine.MAC); err != nil {
			respondError(newHTTPError(400, "Invalid MAC address"), w, r)
			return
		}
		// Because name used in HTTP routes
		if url.PathEscape(machine.Name) != machine.Name {
			respondError(newHTTPError(400, "Invalid name"), w, r)
			return
		}
		machine.Requests = 0
		machine.LastHeartbeat = time.Time{}
		machine.LastRequest = time.Time{}

		err = app.db.Update(func(tx *bolt.Tx) error {
			machines := tx.Bucket(machineBucket)
			data, err := json.Marshal(machine)
			if err != nil {
				return err
			}
			return machines.Put([]byte(machine.Name), []byte(data))
		})
		if err != nil {
			respondError(err, w, r)
			return
		}

		log.Printf("Added machine: %v", machine)
		w.Write([]byte("Ok"))
	})

	r.Post("/remove/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		err := app.db.Update(func(tx *bolt.Tx) error {
			machines := tx.Bucket(machineBucket)
			return machines.Delete([]byte(name))
		})
		if err != nil {
			respondError(err, w, r)
			return
		}
		w.Write([]byte("Deleted record for " + name))
	})

	r.Post("/start/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		// TODO: add heartbeat

		err := app.db.Update(func(tx *bolt.Tx) error {
			machine, err := getMachine(tx, name)
			if err != nil {
				return err
			}
			if machine == nil {
				return newHTTPError(404, "Machine not found")
			}

			machine.Requests++
			machine.LastRequest = time.Now()
			err = startMachine(machine.MAC)
			if err != nil {
				log.Printf("Can't start machine %s: %v", name, err)
			}
			return putMachine(tx, machine)
		})

		if err != nil {
			respondError(err, w, r)
			return
		}

		w.Write([]byte("Accepted"))
	})

	r.Post("/stop/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		err := app.db.Update(func(tx *bolt.Tx) error {
			machine, err := getMachine(tx, name)
			if err != nil {
				return err
			}
			if machine == nil {
				return newHTTPError(404, "Machine not found")
			}

			if machine.Requests > 0 {
				machine.Requests--
			}

			return putMachine(tx, machine)
		})

		if err != nil {
			respondError(err, w, r)
			return
		}

		w.Write([]byte("Accepted"))
	})

	r.Get("/status/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		type resultRecord struct {
			Name      string `json:"name"`
			MAC       string `json:"mac"`
			Requests  int    `json:"requests"`
			IsRunning bool   `json:"isRunning"`
		}
		var result []resultRecord
		err := app.db.View(func(tx *bolt.Tx) error {
			machines := tx.Bucket(machineBucket)
			c := machines.Cursor()
			var err error

			for k, v := c.First(); k != nil && err == nil; k, v = c.Next() {
				machine := machineRecord{}
				err = json.Unmarshal(v, &machine)

				inactivityDuration :=
					time.Duration(config.MachineInactivityTimeoutSeconds) * time.Second
				record := resultRecord{
					machine.Name,
					machine.MAC,
					machine.Requests,
					time.Now().Sub(machine.LastHeartbeat) < inactivityDuration,
				}
				result = append(result, record)
			}
			return err
		})

		if err != nil {
			respondError(err, w, r)
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("/list: %v", err)
		}
	})

	// TODO add auth
	r.Post("/agent/{name}/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		log.Printf("/agent/%s/heartbeat", name)
		response := struct {
			ShouldShutdown bool `json:"shouldShutdown"`
		}{false}

		err := app.db.Batch(func(tx *bolt.Tx) error {
			machine, err := getMachine(tx, name)
			if err != nil {
				return err
			}
			if machine == nil {
				return newHTTPError(404, "Machine not found")
			}

			shutdownDelay := time.Duration(config.ShutdownDelaySeconds) * time.Second
			lastRequestWasLongAgo := time.Now().Sub(machine.LastRequest) > shutdownDelay
			if machine.Requests == 0 && lastRequestWasLongAgo {
				response.ShouldShutdown = true
			}

			machine.LastHeartbeat = time.Now()
			return putMachine(tx, machine)
		})

		if err != nil {
			respondError(err, w, r)
			return
		}

		err = json.NewEncoder(w).Encode(&response)
		if err != nil {
			log.Printf("/agent/%s/heartbeat: %v", name, err)
		}
	})

	log.Printf("Starting at %s", config.BindAddress)
	http.ListenAndServe(config.BindAddress, r)
}

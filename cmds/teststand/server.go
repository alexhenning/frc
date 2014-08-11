package main

import (
	"fmt"
	"log"
	"net/http"

	"go.thoriumrobotics.com/frc.v0/driverstation"
)

// Global mutex
var mutex Mutex

// Global driverstation
var ds *driverstation.DS

func serve(addr string) {
	http.HandleFunc("/take", takeHandler)
	http.HandleFunc("/give", giveHandler)
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/ds", dsHandler)

	ds = driverstation.New(190) // TODO: make team a variable?
	err := ds.Connect()
	if err != nil {
		log.Fatal(err)
	}
	go ds.Run()

	http.ListenAndServe(addr, nil)
}

func takeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	log.Printf("Take from %s (holder=%s)\n", name, mutex.holder)
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Need to specify a name to take the mutex.\n")))
		return
	}

	err := mutex.Lock(name, updateWriter(w))
	if err != nil {
		log.Printf("Error for %s take: %s\n", name, err)
	}
}

func giveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	log.Printf("Give from %s (holder=%s)\n", name, mutex.holder)
	err := mutex.Unlock(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("You don't have permission to give the mutex.\n")))
		return
	}

}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	name, msg := r.FormValue("name"), r.FormValue("message")
	log.Printf("Update from %s. (holder=%s) Message: %s\n", name, mutex.holder, msg)
	mutex.Message(fmt.Sprintf("%s: %s", name, msg))
}

func dsHandler(w http.ResponseWriter, r *http.Request) {
	name, state := r.FormValue("name"), r.FormValue("state")
	log.Printf("DS from %s (holder=%s) State: %s\n", name, mutex.holder, state)
	if name != mutex.holder {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Need to hold the mutex to change the driverstation.\n")))
		return
	}

	if state == "enable" {
		ds.SetEnabled(true)
		mutex.Message("DriverStation Enabled")
	} else if state == "disable" {
		ds.SetEnabled(false)
		mutex.Message("DriverStation Disabled")
	} else if state == "resync" {
		ds.Resync()
		mutex.Message("DriverStation Disabled")
	}
}

func updateWriter(w http.ResponseWriter) Updater {
	return func(s string) error {
		_, err := w.Write([]byte(s + "\n"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return err
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

var mu = sync.RWMutex{}

var kvStore = make(map[string]string)

func main() {

	http.HandleFunc("GET /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /kv/{key} Request received")

		key := r.PathValue("key")

		mu.RLock()
		defer mu.RUnlock()

		value, ok := kvStore[key]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Key not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, value)

	})

	http.HandleFunc("DELETE /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("DELETE /kv/{key} Request received")

		key := r.PathValue("key")

		mu.Lock()
		defer mu.Unlock()

		_, ok := kvStore[key]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "Key not found")
			return
		}

		delete(kvStore, key)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Key deleted")

	})

	http.HandleFunc("PUT /kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("PUT /kv/{key} Request received")

		key := r.PathValue("key")
		defer r.Body.Close()
		value, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error reading body")
			return
		}

		mu.Lock()
		defer mu.Unlock()

		kvStore[key] = string(value)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Key updated")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

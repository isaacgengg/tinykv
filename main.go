package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
	wal  *os.File
}

func NewStore(path string) *Store {
	var err error

	s := Store{data: make(map[string]string)}
	s.wal, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		log.Fatal(err)
	}

	if err := s.replayLog(); err != nil {
		log.Fatal(err)
	}

	return &s
}

func (s *Store) appendLog(op, key, value string) error {
	if op == "PUT" {
		_, err := fmt.Fprintf(s.wal, "PUT %s %s\n", key, value)
		return err
	} else if op == "DELETE" {
		_, err := fmt.Fprintf(s.wal, "DELETE %s\n", key)
		return err
	}
	return fmt.Errorf("unknown op: %s", op)

}

func (s *Store) replayLog() error {
	s.wal.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(s.wal)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Fields(line)

		if len(splitLine) < 2 {
			continue
		}

		method := splitLine[0]
		key := splitLine[1]

		if method == "PUT" {
			if len(splitLine) < 3 {
				continue
			}
			value := splitLine[2]
			s.data[key] = value
		}

		if method == "DELETE" {
			delete(s.data, key)
		}

	}
	return scanner.Err()

}

func (s *Store) getValue(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /kv/{key} Request received")

	key := r.PathValue("key")

	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Key not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, value)
}

func (s *Store) deleteKey(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /kv/{key} Request received")

	key := r.PathValue("key")

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Key not found")
		return
	}

	err := s.appendLog("DELETE", key, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error writing wal.log")
		return
	}
	delete(s.data, key)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Key deleted")

}

func (s *Store) putKeyValue(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /kv/{key} Request received")

	key := r.PathValue("key")
	defer r.Body.Close()
	value, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading body")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	err = s.appendLog("PUT", key, string(value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error writing wal.log")
		return
	}
	s.data[key] = string(value)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Key updated")
}

func main() {

	store := NewStore("wal.log")

	http.HandleFunc("GET /kv/{key}", store.getValue)

	http.HandleFunc("DELETE /kv/{key}", store.deleteKey)

	http.HandleFunc("PUT /kv/{key}", store.putKeyValue)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

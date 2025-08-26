package models


package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Store struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewStore() *Store { return &Store{m: make(map[string]string)} }

func (s *Store) Get(k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[k]
	return v, ok
}
func (s *Store) Set(k, v string) {
	s.mu.Lock()
	s.m[k] = v
	s.mu.Unlock()
}
func (s *Store) Delete(k string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[k]; ok {
		delete(s.m, k)
		return true
	}
	return false
}

func main() {
	tcpPort := flag.String("tcp", "", "TCP port to listen on (e.g. 4000)")
	httpAddr := flag.String("http", "", "HTTP listen address (e.g. :8080)")
	flag.Parse()

	store := NewStore()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// Optional HTTP API (very small: /get?key=, /set?key=&value=, /delete?key=)
	if *httpAddr != "" {
		srv := &http.Server{
			Addr: *httpAddr,
		}
		http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(w, "missing key", http.StatusBadRequest)
				return
			}
			if v, ok := store.Get(key); ok {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, v)
			} else {
				http.Error(w, "not found", http.StatusNotFound)
			}
		})
		http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			val := r.URL.Query().Get("value")
			if key == "" {
				http.Error(w, "missing key", http.StatusBadRequest)
				return
			}
			store.Set(key, val)
			w.WriteHeader(http.StatusOK)
		})
		http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(w, "missing key", http.StatusBadRequest)
				return
			}
			if store.Delete(key) {
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, "not found", http.StatusNotFound)
			}
		})

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("HTTP listening on %s", *httpAddr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP server error: %v", err)
			}
		}()

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(shutCtx)
		}()
	}

	// Optional TCP API
	if *tcpPort != "" {
		ln, err := net.Listen("tcp", ":"+*tcpPort)
		if err != nil {
			log.Fatalf("TCP listen error: %v", err)
		}
		log.Printf("TCP listening on :%s", *tcpPort)

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer ln.Close()
			for {
				conn, err := ln.Accept()
				if err != nil {
					// Likely listener closed during shutdown
					if ctx.Err() != nil {
						return
					}
					log.Printf("Accept error: %v", err)
					continue
				}
				go handleConn(conn, store)
			}
		}()

		go func() {
			<-ctx.Done()
			_ = ln.Close() // unblocks Accept
		}()
	}

	// If neither interface is requested, exit politely.
	if *tcpPort == "" && *httpAddr == "" {
		log.Println("No -tcp or -http provided; nothing to do.")
		return
	}

	<-ctx.Done()
	log.Println("Shutting down…")
	wg.Wait()
	log.Println("Bye.")
}

func handleConn(conn net.Conn, store *Store) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r\n")
		id, result, value := processLine(line, store)
		if _, err := fmt.Fprintf(conn, "%s|%s|%s\n", id, result, value); err != nil {
			return
		}
	}
}

func processLine(line string, store *Store) (id, result, value string) {
	parts := strings.Split(line, "|")
	// Ensure at least id + op + key
	if len(parts) < 3 {
		return safeIdx(parts, 0), "ERROR", ""
	}
	id = parts[0]
	op := strings.ToUpper(parts[1])
	key := parts[2]

	switch op {
	case "GET":
		if key == "" {
			return id, "ERROR", ""
		}
		if v, ok := store.Get(key); ok {
			return id, "OK", v
		}
		return id, "NOTFOUND", ""
	case "SET":
		if len(parts) < 4 {
			return id, "ERROR", ""
		}
		store.Set(key, parts[3])
		return id, "OK", ""
	case "DELETE":
		if key == "" {
			return id, "ERROR", ""
		}
		if store.Delete(key) {
			return id, "OK", ""
		}
		return id, "NOTFOUND", ""
	default:
		return id, "ERROR", ""
	}
}

func safeIdx(parts []string, i int) string {
	if i < len(parts) {
		return parts[i]
	}
	return ""
}

// Got it — here’s what the task is asking, in plain Go terms, plus a small blueprint you can start from.

// What you’re building

// A key–value store server that can speak over TCP (a simple line-based protocol) and optionally over HTTP if the user passes an -http flag.
// The TCP server stays connected and processes multiple newline-terminated requests per connection.

// Command-line flags

// -tcp <PORT> → start the TCP API on this port (e.g., -tcp 4000). If omitted, no TCP listener.

// -http <ADDR> → start the HTTP API on this address (e.g., -http :8080). If omitted, skip HTTP entirely.

// Shared data model

// A concurrent in-memory KV store. (Use sync.Map or a map[string]string guarded by a sync.RWMutex.)

// TCP text protocol

// Request line format (newline terminated):
// REQUEST_ID|OPERATION|KEY|VALUE\n

// REQUEST_ID → opaque string you echo back.

// OPERATION → one of GET, SET, DELETE.

// KEY → string key.

// VALUE → only for SET. For GET/DELETE this field may be empty or omitted (see notes below).

// Response line format:
// REQUEST_ID|RESULT|VALUE\n

// RESULT ∈ OK, NOTFOUND, ERROR

// VALUE → present only when returning a value (i.e., GET that found the key). For OK on SET/DELETE, value is empty.

// Connection semantics: keep the TCP connection open; process requests line-by-line until the client closes it or the server shuts down.

// Parsing notes / errors

// Split on |, trim the trailing \n.

// Minimal validation:

// GET requires at least 3 fields (id, op, key).

// SET requires 4 fields (value required).

// DELETE requires at least 3 fields.

// On malformed input or unknown op: REQUEST_ID|ERROR| (empty value).

// Result rules by operation

// SET key value → store/overwrite; respond id|OK|

// GET key

// if present → id|OK|<value>

// if missing → id|NOTFOUND|

// DELETE key

// if present → delete and id|OK|

// if missing → id|NOTFOUND|

// Concurrency model

// One goroutine per accepted TCP connection.

// Inside each connection, scan lines with bufio.Scanner and handle sequentially.

// Store is shared across connections; protect it for concurrent access.

// Graceful shutdown

// Trap SIGINT/SIGTERM, close listeners, and let active handlers finish.

// For TCP, closing the listener unblocks Accept. For HTTP, call Server.Shutdown(ctx).

// Example I/O

// Client sends (four lines, same connection):

package models

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
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

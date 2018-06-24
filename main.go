package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	srv := NewServer()

	go func() {
		if err := srv.s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-stop

	srv.Shutdown()
}

type Server struct {
	s *http.Server

	// once is used to protect the shutdown code
	once sync.Once
}

// NewServer constructs and returns a Server instance.
func NewServer() *Server {
	mux := http.NewServeMux()

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	srv := &Server{
		s:    s,
		once: sync.Once{},
	}

	mux.HandleFunc("/hash", serveHash)
	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		srv.serveStop(w, r)
	})

	return srv
}

// Shutdown starts the shutdown of the server waiting for existing connections
// to close. Will force shutdown after 5 seconds.
func (srv *Server) Shutdown() {
	srv.once.Do(func() {
		log.Println("Shutting down...")
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		err := srv.s.Shutdown(ctx)
		if err != nil {
			log.Printf("Unable to cleanly shutdown: %v\n", err.Error())
		}
		log.Println("Shutdown complete.")
	})
}

func (srv *Server) serveStop(w http.ResponseWriter, r *http.Request) {
	go srv.Shutdown()
}

func serveHash(w http.ResponseWriter, r *http.Request) {
	wait := time.After(time.Second * 5)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pass := r.PostFormValue("password")
	if pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing 'password' field"))
		return
	}

	encoded := []byte(HashAndEncode(pass))

	<-wait
	w.Write(encoded)
}

// HashAndEncode will return the base64 encoded hash of the given password.Hash
func HashAndEncode(pass string) string {
	hasher := sha512.New()

	hasher.Write([]byte(pass))

	hash := hasher.Sum([]byte{})

	return base64.StdEncoding.EncodeToString(hash)
}

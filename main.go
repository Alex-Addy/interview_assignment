package main

import (
	"crypto/sha512"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/hash", serveHash)

	log.Fatal(http.ListenAndServe(":8080", nil))
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

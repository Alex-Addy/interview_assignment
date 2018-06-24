package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const hashInput = "angryMonkey"
const hashOutput = `ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==`

func TestHash(t *testing.T) {
	output := HashAndEncode(hashInput)

	if output != hashOutput {
		t.Errorf("Expected: %s\nReceived: %s", hashOutput, output)
	}
}

func TestHashEndpoint(t *testing.T) {
	// test happy path
	form := url.Values{}
	form.Set("password", hashInput)
	// TODO figure out why this post data doesn't show up on the other side
	req := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(form.Encode()))
	resp := httptest.NewRecorder()
	serveHash(resp, req)

	// test missing form value
	req = httptest.NewRequest(http.MethodPost, "/hash", nil)
	resp = httptest.NewRecorder()
	serveHash(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Response was not 400 and was instead %d", resp.Code)
	}

	// test incorrect method
	req = httptest.NewRequest(http.MethodGet, "/hash", nil)
	resp = httptest.NewRecorder()
	serveHash(resp, req)
	if resp.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected %d, received %d", http.StatusMethodNotAllowed, resp.Code)
	}
}

func TestShutdown(t *testing.T) {
	srv := NewServer()

	go func() {
		// TODO Add ability to test without binding to localhost:8080
		err := srv.s.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Errorf("Expected ErrServerClosed, got %v", err)
			t.FailNow()
		}
	}()

	// hash request will close this after it is complete
	stop := make(chan interface{}, 0)
	go func() {
		form := url.Values{
			"password": {hashInput},
		}

		resp, err := http.PostForm("http://localhost:8080/hash", form)
		if err != nil {
			t.Fatalf("Request failed with error %v", err)
		}

		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		hash := string(respBody)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d with response %s",
				resp.StatusCode, hash)
		}

		if hash != hashOutput {
			t.Errorf("Response does not match expected hash output.\nExpected: %s\nReceived: %s", hashOutput, hash)
		}

		close(stop)
	}()

	_, err := http.Get("http://localhost:8080/stop")
	if err != nil {
		t.Error(err)
	}

	for range stop {
		log.Println("Got message")
	}
}

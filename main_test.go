package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	form := url.Values{
		"password": {hashInput},
	}
	body := bytes.NewBufferString(form.Encode())
	req := httptest.NewRequest(http.MethodPost, "/hash", body)
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

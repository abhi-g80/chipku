//main_test.go

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	// In case there is an error in forming the request, we fail and stop the test
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(defaultHandler)
	hf.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "chipku v0.1.0"
	actual := recorder.Body.String()
	if string(actual) != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestRandStringBytes(t *testing.T) {
	var expected int = 5
	actual := RandStringBytes(expected)
	if len(actual) != expected {
		t.Errorf("got %v bytes, expected %v,", actual, expected)
	}
}

func TestPastePostHandler(t *testing.T) {
	var str = []byte("hello")
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(str))
	req.Header.Set("Content-Type", "application/json")

	// req, err := http.NewRequest("POST", "", nil)
	// In case there is an error in forming the request, we fail and stop the test
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(pastePostHandler)
	hf.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}

func TestPastePutHandler(t *testing.T) {
	var str = []byte("hello")
	req, err := http.NewRequest("PUT", "/paste", bytes.NewBuffer(str))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("No-Html", "true")

	// req, err := http.NewRequest("POST", "", nil)
	// In case there is an error in forming the request, we fail and stop the test
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(pastePutHandler)
	hf.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

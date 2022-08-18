//main_test.go

package server

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
	hf := http.HandlerFunc(DefaultHandler)
	hf.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "chipku v" + version
	actual := recorder.Body.String()
	if string(actual) != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
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
	hf := http.HandlerFunc(PastePostHandler)
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

	// req, err := http.NewRequest("POST", "", nil)
	// In case there is an error in forming the request, we fail and stop the test
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(PastePutHandler)
	hf.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

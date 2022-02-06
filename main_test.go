//main_test.go

package main

import (
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
	expected := "chipku v0.0.2"
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

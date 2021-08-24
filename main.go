package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
    "syscall"

	"github.com/gorilla/mux"
)

const (
	version string = "0.0.1"
	port    string = ":8080"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	// Attach handlers
    fileServer := http.FileServer(http.Dir("./static"))
    r.Handle("/", fileServer)
	r.HandleFunc("/paste", pasteHandler).Methods("GET")
	r.HandleFunc("/default", defaultHandler).Methods("GET")

	return r
}

func main() {
	l := log.New(os.Stdout, "chipku -> ", log.LstdFlags)
	r := newRouter()

	s := &http.Server{
		Addr:         port,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		l.Println("Starting server on port", port)

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
    signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	sig := <-sigChan

	l.Printf("Received %s, gracefully shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()
	s.Shutdown(tc)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "chipku v%s", version)
}

func pasteHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }
    fmt.Fprintf(w, "POST request successful\n")
    paste := r.FormValue("paste")
    fmt.Fprintf(w, "paste = %s\n", paste)
}

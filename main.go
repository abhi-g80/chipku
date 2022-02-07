package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const (
	version string = "0.0.2"
	port    string = ":8080"
)

var Chipkus = map[string]string{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	// Attach handlers
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/", fileServer)
	r.HandleFunc("/paste", pasteHandler).Methods("POST")
	r.HandleFunc("/default", defaultHandler).Methods("GET")
	r.HandleFunc("/{id}", fetchHandler).Methods("GET")

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

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		fmt.Printf("Error while fetching vars")
		return
	}
	if x, found := Chipkus[id]; found {
		fmt.Fprintf(w, "%s", x)
	} else {
		fmt.Fprintf(w, "Invalid id (%s) provided :(", id)
	}
}

var letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func store(v string) string {
	hash_val := RandStringBytes(6)
	Chipkus[hash_val] = v
	return hash_val
}

func pasteHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	value := r.FormValue("paste-area")
	hash_val := store(value)
	fmt.Fprintf(w, "%s\n", hash_val)
}

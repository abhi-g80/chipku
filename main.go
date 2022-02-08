package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

// Global logger (bad !) think of middleware
var logger = log.New(os.Stdout, "chipku -> ", log.LstdFlags|log.Lmicroseconds)

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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

    logger.Println("attaching handlers")

	// Attach handlers
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/", fileServer)
	r.HandleFunc("/paste", pastePostHandler).Methods("POST")
	r.HandleFunc("/paste", pastePutHandler).Methods("PUT")
	r.HandleFunc("/default", defaultHandler).Methods("GET")
	r.HandleFunc("/{id}", fetchHandler).Methods("GET")

    logger.Println("return router object")

	return r
}

func main() {
    logger.Printf("starting chipku v%s", version)
	r := newRouter()

	s := &http.Server{
		Addr:         port,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		logger.Println("starting server on port", port)

		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	sig := <-sigChan

	logger.Printf("received %s, gracefully shutdown", sig)

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
		logger.Println("Error while fetching vars")
		return
	}
	if x, found := Chipkus[id]; found {
		fmt.Fprintf(w, "%s", x)
	} else {
		fmt.Fprintf(w, "Invalid id (%s) provided :(", id)
	}
}

func pastePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	value := r.FormValue("paste-area")
	hash_val := store(value)
	fmt.Fprintf(w, "%s", hash_val)
}

func pastePutHandler(w http.ResponseWriter, r *http.Request) {
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        logger.Printf("error while reading body")
    }
    hash_val := store(string(b))
    fmt.Fprintf(w, "%s", hash_val)
}

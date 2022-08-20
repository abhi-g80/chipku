package server

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	InfoLog.Printf("attaching handlers")

	r.Handle("/", IndexFileServer())
	r.HandleFunc("/paste", PastePostHandler).Methods("POST")
	r.HandleFunc("/paste", PastePutHandler).Methods("PUT")
	r.HandleFunc("/default", DefaultHandler).Methods("GET")
	r.HandleFunc("/{hashID}", FetchHandler).Methods("GET")

	return r
}

// Serve exported
func Serve(port string) {
	InfoLog.Printf("starting chipku v%s", Version)
	r := newRouter()

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		InfoLog.Printf("using port %s", port)

		err := s.ListenAndServe()
		if err != nil {
			ErrorLog.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan

	DebugLog.Printf("received %s, gracefully shutting down", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err := s.Shutdown(tc)
	if err != nil {
		ErrorLog.Printf("could not shutdown gracefully %s", err)
	}
}

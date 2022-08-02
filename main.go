package main

import (
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const (
	version string = "0.0.3"
	port    string = ":8080"
)

var Chipkus = map[string]string{}

// Global logger (bad !) think of middleware
var logger = log.New(os.Stdout, "[\033[0;34mchipku\033[0m] ", log.LstdFlags|log.Lmicroseconds)

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

    logger.Println("Info  -> attaching handlers")

	// Attach handlers
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/", fileServer)
	r.HandleFunc("/paste", pastePostHandler).Methods("POST")
	r.HandleFunc("/paste", pastePutHandler).Methods("PUT")
	r.HandleFunc("/default", defaultHandler).Methods("GET")
	r.HandleFunc("/{id}", fetchHandler).Methods("GET")

    logger.Println("Info  -> return router object")

	return r
}

func main() {
    logger.Printf("Info  -> starting chipku v%s", version)
	r := newRouter()

	s := &http.Server{
		Addr:         port,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		logger.Println("Info  -> starting server on port", port)

		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	sig := <-sigChan

	logger.Printf("Info  -> received %s, gracefully shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()
	s.Shutdown(tc)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "chipku v%s", version)
}

type CodeData struct {
    Code string
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		logger.Println("Error -> while fetching vars")
		return
	}
	if x, found := Chipkus[id]; found {
        ts, err := template.ParseFiles("./static/code.html.tmpl")
        if err != nil {
            logger.Println("Error -> while loading code template")
            logger.Printf("Error -> %s", err)
            return
        }
        var y []string
        for _, line := range strings.Split(strings.TrimSuffix(x, "\n"), "\n") {
            tmp := "<code>" + html.EscapeString(line) + "</code>"
            y = append(y, tmp)
        }
        z := strings.Join(y, "")
        data := CodeData{Code: z}
        err = ts.Execute(w, data)
        if err != nil {
            logger.Println("Error -> something went wrong while templating code")
            logger.Printf("Error -> %s", err)
        }
		// fmt.Fprintf(w, "%s", x)
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
        logger.Printf("Error -> while reading body")
    }
    hash_val := store(string(b))
    fmt.Fprintf(w, "%s", hash_val)
}

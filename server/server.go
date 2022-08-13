package server

import (
	"context"
	"embed"
	"fmt"
	"html"
	"io/fs"
	"io/ioutil"
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

//go:embed static/index.html
var index embed.FS

//go:embed static/code.html.tmpl
var code_template embed.FS

var Chipkus = map[string]string{}

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

	LogInfo("attaching handlers")

	fsys := fs.FS(index)
	html, _ := fs.Sub(fsys, "static")

	fileServer := http.FileServer(http.FS(html))
	r.Handle("/", fileServer)
	r.HandleFunc("/paste", pastePostHandler).Methods("POST")
	r.HandleFunc("/paste", pastePutHandler).Methods("PUT")
	r.HandleFunc("/default", defaultHandler).Methods("GET")
	r.HandleFunc("/{hash_id}", fetchHandler).Methods("GET")

	return r
}

func usage() {
	fmt.Println("Usage: chipku <port>")
	fmt.Println("default port is 8080")
	os.Exit(0)
}

func Serve(port string) {
	// var port string = ":8080"

	// if len(os.Args) == 2 {
	// 	port = ":" + os.Args[1]
	// } else if len(os.Args) > 2 {
	// 	usage()
	// }

	LogInfo("starting chipku v%s", Version)
	r := newRouter()

	s := &http.Server{
        Addr:         ":"+port,
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		LogInfo("using port %s", port)

		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	sig := <-sigChan

	LogDebug("received %s, gracefully shutting down", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()
	s.Shutdown(tc)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "chipku v%s", Version)
}

type CodeData struct {
	Code string
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash_id, ok := vars["hash_id"]
	if !ok {
		LogError("something went wrong while fetching vars %v", vars)
		return
	}
	split := strings.Split(hash_id, ".")
	id := split[0]
	var lang string = "plaintext"
	if len(split) > 1 {
		lang = split[1]
	}
	if x, found := Chipkus[id]; found {
		v, _ := r.Header["No-Html"]
		if v != nil {
			w.Header().Add("Content-Type", "text; charset=UTF-8")
			fmt.Fprintf(w, "%s", x)
			return
		}
		fsys := fs.FS(code_template)
		ts, err := template.ParseFS(fsys, "static/code.html.tmpl")
		if err != nil {
			LogError("could not load code template 😔")
			LogError("%s", err)
			return
		}
		var y []string
		for _, line := range strings.Split(strings.TrimSuffix(x, "\n"), "\n") {
			tmp := "<code class=\"language-" + lang + "\">" + html.EscapeString(line) + "</code>"
			y = append(y, tmp)
		}
		z := strings.Join(y, "")
		data := CodeData{Code: z}
		err = ts.Execute(w, data)
		if err != nil {
			LogError("something went wrong while templating code %s", err)
		}
	} else {
		fmt.Fprintf(w, "Invalid id %s provided :(", id)
	}
}

func pastePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	value := r.FormValue("paste-area")
	hash_val := store(value)
	url := "/" + hash_val
	LogInfo("new %s request from connection from %s", r.Method, r.RemoteAddr)
	LogInfo("User-agent %s", r.UserAgent())
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func pastePutHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Printf("\033[0;31m[Error]\033[0m -> while reading body")
	}
	hash_val := store(string(b))
	LogInfo("new %s request from connection from %s", r.Method, r.RemoteAddr)
	LogInfo("User-agent %s", r.UserAgent())
	fmt.Fprintf(w, "%s", hash_val)
}

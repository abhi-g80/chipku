package server

import (
	"embed"
	"fmt"
	"html"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

//go:embed static/index.html
var index embed.FS

//go:embed static/code.html.tmpl
var codeTemplate embed.FS

var chipkus = map[string]string{}

var letterBytes = "abcdefghijklmnopqrstuvwxyz"

// CodeData exported
type CodeData struct {
	Code string
}

// RandStringBytes return a random string bytes of size n
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// IndexFileServer read the static html and return the http handler
func IndexFileServer() http.Handler {
	fsys := fs.FS(index)
	html, _ := fs.Sub(fsys, "static")

	return http.FileServer(http.FS(html))
}

func store(v string) string {
	hashVal := RandStringBytes(6)
	chipkus[hashVal] = v
	return hashVal
}

// DefaultHandler return the chipku version
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "chipku v%s", Version)
}

func getCodeTemplate(templ embed.FS) *template.Template {
	fsys := fs.FS(templ)
	ts, err := template.ParseFS(fsys, "static/code.html.tmpl")
	if err != nil {
		LogError("could not load code template 😔")
		LogError("%s", err)
	}
	return ts
}

func enrichWithHTMLTags(codedata string, language string) CodeData {
	var y []string
	for _, line := range strings.Split(strings.TrimSuffix(codedata, "\n"), "\n") {
		tmp := "<code class=\"language-" + language + "\">" + html.EscapeString(line) + "</code>"
		y = append(y, tmp)
	}
	z := strings.Join(y, "")
	return CodeData{Code: z}
}

// FetchHandler GET handler function for fetching the code snippets
func FetchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashID, ok := vars["hashID"]
	if !ok {
		LogError("something went wrong while fetching vars %v", vars)
		return
	}
	split := strings.Split(hashID, ".")
	id := split[0]
	var lang string = "plaintext"
	if len(split) > 1 {
		lang = split[1]
	}
	if data, found := chipkus[id]; found {
		_, ok := r.Header["No-Html"]
		if ok {
			w.Header().Add("Content-Type", "text; charset=UTF-8")
			fmt.Fprintf(w, "%s", data)
			return
		}
		ts := getCodeTemplate(codeTemplate)
		enrichedData := enrichWithHTMLTags(data, lang)
		err := ts.Execute(w, enrichedData)
		if err != nil {
			LogError("something went wrong while templating code %s", err)
		}
	} else {
		LogInfo("invalid id %s requested by %s", id, r.RemoteAddr)
		fmt.Fprintf(w, "Invalid id %s provided :(", id)
	}
}

// PastePostHandler POST handler func for managing snippets being posted via form
func PastePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		LogError("ParseForm() err: %v", err)
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	value := r.FormValue("paste-area")
	hashVal := store(value)
	url := "/" + hashVal
	LogInfo("new %s request from connection from %s", r.Method, r.RemoteAddr)
	LogInfo("User-agent %s", r.UserAgent())
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// PastePutHandler PUT handler func for managing snippets sent via PUT
func PastePutHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogError("while reading body = %v", b)
		return
	}
	hashVal := store(string(b))
	LogInfo("new %s request from connection from %s", r.Method, r.RemoteAddr)
	LogInfo("User-agent %s", r.UserAgent())
	fmt.Fprintf(w, "%s", hashVal)
}

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

	"github.com/labstack/echo/v4"
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

func store(v string) string {
	hashVal := RandStringBytes(6)
	chipkus[hashVal] = v
	return hashVal
}

// IndexFileServer read the static html and return the http handler
func IndexFileServer() http.Handler {
	fsys := fs.FS(index)
	html, _ := fs.Sub(fsys, "static")

	return http.FileServer(http.FS(html))
}

func version(c echo.Context) error {
	return c.String(http.StatusOK, "Chipku v"+Version)
}

func getCodeTemplate(templ embed.FS) (*template.Template, error) {
	fsys := fs.FS(templ)
	ts, err := template.ParseFS(fsys, "static/code.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("could not template: %v", err)
	}
	return ts, nil
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

func fetchHandler(c echo.Context) error {
	hashVal := c.Param("hashVal")
	var lang string = "plaintext"
	split := strings.Split(hashVal, ".")
	id := split[0]
	if len(split) > 1 {
		lang = split[1]
	}
	data, found := chipkus[id]
	if !found {
		return c.String(http.StatusNotFound, "requested hash not found")
	}
	_, ok := c.Request().Header["X-No-Html"]
	if ok {
		c.Response().Header().Add("Content-Type", "text; charset=UTF-8")
		return c.String(http.StatusOK, data)
	}
	ts, err := getCodeTemplate(codeTemplate)
	if err != nil {
		return err
	}
	enrichedData := enrichWithHTMLTags(data, lang)
	return ts.Execute(c.Response().Writer, enrichedData)
}

func pastePostHandler(c echo.Context) error {
	value := c.FormValue("paste-area")
	hashVal := store(value)
	url := "/" + hashVal
	return c.Redirect(http.StatusSeeOther, url)
}

func pastePutHandler(c echo.Context) error {
	value, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("while reading body: %v", err)
	}
	hashVal := store(string(value))
	c.Logger().Infof("responding %s", hashVal)
	return c.String(http.StatusOK, hashVal)
}

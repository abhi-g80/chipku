package server

import (
	"fmt"
	"html"
	"strings"
	"testing"
)

func TestRandStringBytes(t *testing.T) {
	var expected int = 6
	actual := RandStringBytes(expected)
	if len(actual) != expected {
		t.Errorf("got %v bytes, expected %v,", actual, expected)
	}
}

func TestEnrichWithHtmlTags(t *testing.T) {
	var lang string = "lol"
	var codedata string = "test1"
	var expectedString string = fmt.Sprintf("<code class=\"language-%s\">%s</code>", lang, codedata)
	var expected CodeData = CodeData{Code: expectedString}
	var actual CodeData = enrichWithHTMLTags(codedata, lang)

	if strings.Compare(actual.Code, expected.Code) != 0 {
		t.Errorf("got %v, expected %v", actual, expected)
	}
}

func TestEnrichWithHtmlTagsWithHtmlstring(t *testing.T) {
	var lang string = "lol"
	var codedata string = "<html><body>help</body></html>"
	var expectedString string = fmt.Sprintf("<code class=\"language-%s\">%s</code>", lang, html.EscapeString(codedata))
	var expected CodeData = CodeData{Code: expectedString}
	var actual CodeData = enrichWithHTMLTags(codedata, lang)

	if strings.Compare(actual.Code, expected.Code) != 0 {
		t.Errorf("got %v, expected %v", actual, expected)
	}
}

package domutil

import (
	"bytes"
	"golang.org/x/net/html"
	"io/ioutil"
	"testing"
)

func TestTokenize(t *testing.T) {
	buf, err := ioutil.ReadFile("testdata/tokenize.html")
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewReader(buf)

	divDepth := 0
	Tokenize(r, func(token html.Token) (bool, error) {
		switch token.Type {
		case html.StartTagToken:
			if token.Data == "div" {
				divDepth++
			}

		case html.TextToken:
			if divDepth == 1 && token.Data != "test" {
				t.Error("div text is not 'test'")
			}

		case html.EndTagToken:
			if token.Data == "div" {
				divDepth--
				return false, nil // stop walking dom tree.
			}
			if token.Data == "body" {
				t.Error("not stopped")
			}
		}

		return true, nil
	})
}

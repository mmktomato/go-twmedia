package domutil

import (
	"golang.org/x/net/html"
	"io"
)

func Tokenize(r io.Reader, fn func(token html.Token) (bool, error)) error {
	tokenizer := html.NewTokenizer(r)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return nil
			}
			return tokenizer.Err()

		default:
			doContinue, err := fn(tokenizer.Token())
			if err != nil {
				return err
			}
			if !doContinue {
				return nil
			}
		}
	}
	return nil
}

func FindAttr(attrs []html.Attribute, key string, fn func(html.Attribute) error) error {
	for _, v := range attrs {
		if v.Key == key {
			return fn(v)
		}
	}
	return nil
}

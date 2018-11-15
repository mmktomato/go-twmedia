package twparser

import (
	"golang.org/x/net/html"
	"io"
	neturl "net/url"
	"strings"
)

type TwMedia struct {
	ImageUrls map[string]string // url : filename
	VideoUrl  string
}

func ParseTweet(r io.Reader) (*TwMedia, error) {
	res := &TwMedia{make(map[string]string), ""}
	tokenizer := html.NewTokenizer(r)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return res, nil
			}
			return nil, tokenizer.Err()

		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" {
				err := parseMetaAttr(token.Attr, res)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return res, nil
}

func parseMetaAttr(attrs []html.Attribute, twMedia *TwMedia) error {
	return findAttr(attrs, "property", func(propAttr html.Attribute) error {
		return findAttr(attrs, "content", func(contentAttr html.Attribute) error {
			switch propAttr.Val {
			case "og:image":
				if isTargetImage(contentAttr.Val) {
					filename, err := getImageFilename(contentAttr.Val)
					if err != nil {
						return err
					}
					twMedia.ImageUrls[contentAttr.Val] = filename
				}
			case "og:video:url":
				twMedia.VideoUrl = contentAttr.Val
			}
			return nil
		})
	})
}

func findAttr(attrs []html.Attribute, key string, fn func(html.Attribute) error) error {
	for _, v := range attrs {
		if v.Key == key {
			return fn(v)
		}
	}
	return nil
}

func isTargetImage(url string) bool {
	return strings.HasSuffix(url, ":large")
}

func getImageFilename(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}

	tokens := strings.Split(u.Path, "/")
	ret := tokens[len(tokens)-1]
	return strings.Split(ret, ":")[0], nil // foo.jpg:large
}

package twparser

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type TwMedia struct {
	ImageUrls []string
	VideoUrl  string
}

func ParseTweet(r io.Reader) (*TwMedia, error) {
	res := &TwMedia{}
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
				parseMetaAttr(token.Attr, res)
			}
		}
	}
	return res, nil
}

func parseMetaAttr(attrs []html.Attribute, twMedia *TwMedia) {
	findAttr(attrs, "property", func(propAttr html.Attribute) {
		findAttr(attrs, "content", func(contentAttr html.Attribute) {
			switch propAttr.Val {
			case "og:image":
				if isTargetImage(contentAttr.Val) {
					twMedia.ImageUrls = append(twMedia.ImageUrls, contentAttr.Val)
				}
			case "og:video:url":
				twMedia.VideoUrl = contentAttr.Val
			}
		})
	})
}

func findAttr(attrs []html.Attribute, key string, fn func(html.Attribute)) {
	for _, v := range attrs {
		if v.Key == key {
			fn(v)
		}
	}
}

func isTargetImage(url string) bool {
	return strings.HasSuffix(url, ":large")
}

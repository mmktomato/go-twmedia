// TODO: separate this file out of 'main' package.

package main

import (
	"golang.org/x/net/html"
	"io"
)

type TwMedia struct {
	imageUrls []string
	videoUrl string
}

func ParseTweet(r io.Reader) (*TwMedia, error) {
	res := &TwMedia{}
	tokenizer := html.NewTokenizer(r)

	LOOP:
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				break LOOP
			}
			return nil, tokenizer.Err()

		case html.StartTagToken: fallthrough
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
				twMedia.imageUrls = append(twMedia.imageUrls, contentAttr.Val)
			case "og:video:url":
				twMedia.videoUrl = contentAttr.Val
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

package twparser

import (
	"golang.org/x/net/html"
	"io"
	neturl "net/url"
	"strings"

	"github.com/mmktomato/go-twmedia/twparser/video"
)

type TwMedia struct {
	ImageUrls map[string]string // url : filename
	VideoUrl  string
}

func ParseTweet(r io.Reader) (*TwMedia, error) {
	res := &TwMedia{make(map[string]string), ""}

	err := tokenize(r, func(token html.Token) error {
		switch token.Type {
		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			if token.Data == "meta" {
				err := parseMetaAttr(token.Attr, res)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return res, err
}

// ParseVideo returns playlist's url. Typically it has `.m3u8` extension.
func ParseVideo(r io.Reader) (ret string, err error) {
	// TODO: write test code.

	err = tokenize(r, func(token html.Token) error {
		switch token.Type {
		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			if token.Data == "script" {
				return findAttr(token.Attr, "src", func(srcAttr html.Attribute) error {
					// TODO move to other package
					if strings.Contains(srcAttr.Val, "TwitterVideoPlayerIframe") {
						jsurl := srcAttr.Val
						token, err := video.GetAuthToken(jsurl)
						if err != nil {
							return err
						}
						ret = token
					}
					return nil
				})
			}
		}
		return nil
	})
	return ret, err
}

func tokenize(r io.Reader, fn func(token html.Token) error) error {
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
			err := fn(tokenizer.Token())
			if err != nil {
				return err
			}
		}
	}
	return nil
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

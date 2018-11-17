package twparser

import (
	"golang.org/x/net/html"
	"io"
	neturl "net/url"
	"strings"

	"github.com/mmktomato/go-twmedia/twparser/domutil"
	"github.com/mmktomato/go-twmedia/twparser/video"
)

type TwMedia struct {
	ImageUrls map[string]string // url : filename
	VideoUrl  string
}

func ParseTweet(r io.Reader) (*TwMedia, error) {
	res := &TwMedia{make(map[string]string), ""}

	err := domutil.Tokenize(r, func(token html.Token) (bool, error) {
		switch token.Type {
		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			if token.Data == "meta" {
				err := parseMetaAttr(token.Attr, res)
				if err != nil {
					return false, err
				}
			}
		}
		return true, nil
	})

	return res, err
}

// ParseVideo returns playlist's url. Typically it has `.m3u8` extension.
func ParseVideo(r io.Reader) (ret string, err error) {
	// TODO: write test code. needs mock for `Fetch` because video.GetAuthToken uses it.

	err = domutil.Tokenize(r, func(token html.Token) (bool, error) {
		switch token.Type {
		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			if token.Data == "script" {
				authToken, err := video.GetAuthToken(token.Attr)
				if err != nil {
					return false, err
				}
				ret = authToken // temporary. TODO: `ret` is m3u8 url.
				if authToken != "" {
					return false, nil
				}
			}
		}
		return true, nil
	})
	return ret, err
}

func parseMetaAttr(attrs []html.Attribute, twMedia *TwMedia) error {
	return domutil.FindAttr(attrs, "property", func(propAttr html.Attribute) error {
		return domutil.FindAttr(attrs, "content", func(contentAttr html.Attribute) error {
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

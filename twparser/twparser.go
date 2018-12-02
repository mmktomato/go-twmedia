package twparser

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	neturl "net/url"
	"regexp"
	"strings"

	"github.com/mmktomato/go-twmedia/twparser/domutil"
)

type TwMedia struct {
	ImageUrls map[string]string // url : filename
	VideoUrl  string
	TweetId   string
}

var tweetIdRegex = regexp.MustCompile(`^https://twitter.com/[^/]+/status/([^/]+)/?`)

func ParseTweet(tweetUrl string, r io.Reader) (*TwMedia, error) {
	tweetId := getTweetId(tweetUrl)
	if tweetId == "" {
		return nil, errors.New("TweetId not found")
	}

	res := &TwMedia{
		ImageUrls: make(map[string]string),
		VideoUrl:  "",
		TweetId:   tweetId,
	}

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
		case html.EndTagToken:
			if token.Data == "head" {
				return false, nil
			}
		}
		return true, nil
	})

	return res, err
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

func getTweetId(url string) string {
	found := tweetIdRegex.FindStringSubmatch(url)
	if 1 < len(found) {
		return found[1]
	}
	return ""
}

package twparser

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	neturl "net/url"
	"regexp"
	"strings"

	"github.com/mmktomato/go-twmedia/twparser/domutil"
	"github.com/mmktomato/go-twmedia/twparser/video"
)

type TwMedia struct {
	ImageUrls map[string]string // url : filename
	VideoUrl  string
}

var tweetIdRegex = regexp.MustCompile(`^https://twitter.com/[^/]+/status/([^/]+)/?`)

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

// ParseVideo returns *video.TrackInfo. It contains playlist url and content id.
// Typically playlist url has `.m3u8` extension.
func ParseVideo(tweetUrl string, r io.Reader) (ret *video.TrackInfo, err error) {
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
				if authToken != "" {
					tweetId := getTweetId(tweetUrl)
					if tweetId == "" {
						return false, errors.New("TweetId not found")
					}
					ret, err = video.GetTrackInfo(tweetId, authToken)

					return false, err
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

func getTweetId(url string) string {
	found := tweetIdRegex.FindStringSubmatch(url)
	if 1 < len(found) {
		return found[1]
	}
	return ""
}

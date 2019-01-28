package tw

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"regexp"

	"github.com/mmktomato/go-twmedia/svc/image"
	"github.com/mmktomato/go-twmedia/util"
	"github.com/mmktomato/go-twmedia/util/domutil"
)

type Tweet struct {
	ImageUrls map[string]string // url : filename
	VideoUrl  string
	TweetId   string
}

type TweetService interface {
	ParseTweet(string, io.Reader) (*Tweet, error)
}

type TweetServiceImpl struct {
	imageService image.ImageService
	logger       *util.TinyLogger
}

func NewTweetServiceImpl(imageService image.ImageService, logger *util.TinyLogger) *TweetServiceImpl {
	return &TweetServiceImpl{imageService, logger}
}

var tweetIdRegex = regexp.MustCompile(`^https://twitter.com/[^/]+/status/([^/]+)/?`)

func (svc *TweetServiceImpl) ParseTweet(tweetUrl string, r io.Reader) (*Tweet, error) {
	tweetId := svc.getTweetId(tweetUrl)
	if tweetId == "" {
		return nil, errors.New("TweetId not found")
	}

	res := &Tweet{
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
				err := svc.parseMetaAttr(token.Attr, res)
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

func (svc *TweetServiceImpl) parseMetaAttr(attrs []html.Attribute, tweet *Tweet) error {
	return domutil.FindAttr(attrs, "property", func(propAttr html.Attribute) error {
		return domutil.FindAttr(attrs, "content", func(contentAttr html.Attribute) error {
			switch propAttr.Val {
			case "og:image":
				if svc.imageService.IsTargetImage(contentAttr.Val) {
					filename, err := svc.imageService.GetImageFilename(contentAttr.Val)
					if err != nil {
						return err
					}
					tweet.ImageUrls[contentAttr.Val] = filename
					svc.logger.Verbosef("%s -> %s\n", contentAttr.Val, filename)
				} else {
					svc.logger.Verbosef("%s is not a target image\n", contentAttr.Val)
				}
			case "og:video:url":
				tweet.VideoUrl = contentAttr.Val
			}
			return nil
		})
	})
}

func (svc *TweetServiceImpl) getTweetId(url string) string {
	found := tweetIdRegex.FindStringSubmatch(url)
	if 1 < len(found) {
		for i := range found[1:] {
			index := i + 1
			svc.logger.Verbosef("tweetIdRegex found[%d]: %v\n", index, found[index])
		}
		return found[1]
	}
	return ""
}

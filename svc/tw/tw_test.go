package tw

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/mmktomato/go-twmedia/svc/image"
	"github.com/mmktomato/go-twmedia/util"
)

var (
	logger = util.NewTinyLogger(false)
	svc    = NewTweetServiceImpl(image.NewImageServiceImpl(logger), logger)
)

func TestParseTweet(t *testing.T) {
	tests := []struct {
		testfile string
		images   map[string]string
		video    string
		tweetUrl string
		tweetId  string
	}{
		{
			"testdata/tweet_no_media.html",
			make(map[string]string),
			"",
			"https://twitter.com/some_account/status/id1",
			"id1",
		},
		{
			"testdata/tweet_with_image.html",
			map[string]string{
				"https://example.com/image1.jpg:large": "image1.jpg",
				"https://example.com/image2.jpg:large": "image2.jpg",
			},
			"",
			"https://twitter.com/some_account/status/id1",
			"id1",
		},
		{
			"testdata/tweet_with_video.html",
			make(map[string]string),
			"https://example.com/0123456789",
			"https://twitter.com/some_account/status/id1",
			"id1",
		},
	}

	for _, tt := range tests {
		buf, err := ioutil.ReadFile(tt.testfile)
		if err != nil {
			t.Fatal(err)
		}

		r := bytes.NewReader(buf)
		twMedia, err := svc.ParseTweet(tt.tweetUrl, r)
		if err != nil {
			t.Fatal(err)
		}

		// ImageUrls
		if len(twMedia.ImageUrls) != len(tt.images) {
			t.Errorf("%s: length not match", tt.testfile)
		}
		for expectedUrl, expectedFilename := range tt.images {
			filename, ok := twMedia.ImageUrls[expectedUrl]
			if !ok {
				t.Errorf("%s: image not found (%s)", tt.testfile, expectedUrl)
			}
			if filename != expectedFilename {
				t.Errorf("%s: image not match (%s)", tt.testfile, expectedFilename)
			}
		}

		// VideoUrl
		if twMedia.VideoUrl != tt.video {
			t.Errorf("%s: video not match", tt.testfile)
		}

		// TweetId
		if twMedia.TweetId != tt.tweetId {
			t.Errorf("%s: tweetId not match", tt.testfile)
		}
	}
}

func TestGetTweetId(t *testing.T) {
	tests := []struct {
		url, expectedId string
	}{
		{"https://twitter.com/some_account/status/id1", "id1"},
		{"https://twitter.com/some_account/status/id2/photo/1", "id2"},
	}

	for _, tt := range tests {
		res := svc.getTweetId(tt.url)
		if res != tt.expectedId {
			t.Errorf("%s: id not match", tt.url)
		}
	}
}

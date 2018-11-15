package twparser

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestParseTweet(t *testing.T) {
	tests := []struct {
		testfile string
		images   map[string]string
		video    string
	}{
		{"testdata/tweet_no_media.html", make(map[string]string), ""},
		{
			"testdata/tweet_with_image.html",
			map[string]string{
				"https://example.com/image1.jpg:large": "image1.jpg",
				"https://example.com/image2.jpg:large": "image2.jpg",
			},
			"",
		},
		{"testdata/tweet_with_video.html", make(map[string]string), "https://example.com/0123456789"},
	}

	for _, tt := range tests {
		buf, err := ioutil.ReadFile(tt.testfile)
		if err != nil {
			t.Fatal(err)
		}

		r := bytes.NewReader(buf)
		twMedia, err := ParseTweet(r)
		if err != nil {
			t.Fatal(err)
		}

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
		if twMedia.VideoUrl != tt.video {
			t.Errorf("%s: video not match", tt.testfile)
		}
	}
}

func TestGetImageFilename(t *testing.T) {
	tests := []struct {
		url, filename string
		isErr         bool
	}{
		{"http://example.com/image1.jpg:large", "image1.jpg", false},
		{"://example.com/image1.jpg:large", "", true},
	}

	for _, tt := range tests {
		res, err := getImageFilename(tt.url)
		if tt.isErr && err == nil {
			t.Errorf("%s: error not found", tt.url)
		}
		if !tt.isErr && err != nil {
			t.Errorf("%s: error occured", tt.url)
		}
		if tt.filename != res {
			t.Errorf("%s: filename not match", tt.url)
		}
	}
}

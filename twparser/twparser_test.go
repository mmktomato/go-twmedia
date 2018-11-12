package twparser

import (
	"testing"
	"bytes"
	"io/ioutil"
)

func TestParseTweet_noMedia(t *testing.T) {
	tests := []struct {
		file string
		images []string
		video string
	}{
		{ "testdata/tweet_no_media.html", make([]string, 0), "" },
		{ "testdata/tweet_with_image.html", []string{ "https://example.com/image1.jpg:large",  "https://example.com/image2.jpg:large" }, "" },
		{ "testdata/tweet_with_video.html", make([]string, 0), "https://example.com/video.mp4" },
	}

	for _, tt := range tests {
		buf, err := ioutil.ReadFile(tt.file)
		if err != nil {
			t.Fatal(err)
		}

		r := bytes.NewReader(buf)
		twMedia, err := ParseTweet(r)
		if err != nil {
			t.Fatal(err)
		}

		if len(twMedia.imageUrls) != len(tt.images) {
			t.Errorf("%s: length not match", tt.file)
		}
		for i, v := range twMedia.imageUrls {
			if v != tt.images[i] {
				t.Errorf("%s: image not match", tt.file)
			}
		}
		if twMedia.videoUrl != tt.video {
			t.Errorf("%s: video not match", tt.file)
		}
	}
}

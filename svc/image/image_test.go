package image

import (
	"testing"

	"github.com/mmktomato/go-twmedia/util"
)

var svc = NewImageServiceImpl(util.NewTinyLogger(false))

func TestGetImageFilename(t *testing.T) {
	tests := []struct {
		url, filename string
		isErr         bool
	}{
		{"http://example.com/image1.jpg:large", "image1.jpg", false},
		{"://example.com/image1.jpg:large", "", true},
	}

	for _, tt := range tests {
		res, err := svc.GetImageFilename(tt.url)
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

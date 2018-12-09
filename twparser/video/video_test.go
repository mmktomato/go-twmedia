package video

import (
	"testing"

	"github.com/grafov/m3u8"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestFindBiggestVideo(t *testing.T) {
	small := &m3u8.Variant{URI: "/sample/small.m3u8", VariantParams: m3u8.VariantParams{Bandwidth: 10}}
	medium := &m3u8.Variant{URI: "/sample/medium.m3u8", VariantParams: m3u8.VariantParams{Bandwidth: 20}}
	big := &m3u8.Variant{URI: "/sample/big.m3u8", VariantParams: m3u8.VariantParams{Bandwidth: 30}}

	tests := []struct {
		variants []*m3u8.Variant
	}{
		{[]*m3u8.Variant{small, medium, big}},
		{[]*m3u8.Variant{small, big, medium}},
		{[]*m3u8.Variant{big, medium, small}},
	}

	for i, tt := range tests {
		res := findBiggestVideo(tt.variants)
		if res != big {
			t.Errorf("%d: big is not returned.", i)
		}
	}
}

func TestExtractAuthToken(t *testing.T) {
	tests := []struct {
		jsresp, token string
	}{
		{`foobar authorization:"myToken" foobar`, "myToken"},
		{`foobar authorization: "myToken" foobar`, "myToken"},
		{`foobar`, ""},
	}

	for i, tt := range tests {
		httpmock.Activate()
		defer httpmock.Deactivate()

		jsurl := "https://localhost/TwitterVideoPlayerIframe.js"
		httpmock.RegisterResponder("GET", jsurl, httpmock.NewStringResponder(200, tt.jsresp))

		token, err := extractAuthToken(jsurl)
		if token != tt.token {
			t.Errorf("%d: unexpected token -> %s", i, token)
			t.Errorf("%d: err -> %v", i, err)
		}
	}
}

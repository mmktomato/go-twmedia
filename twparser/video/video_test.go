package video

import (
	"testing"

	"github.com/grafov/m3u8"
)

func TestFindBiggestVideo(t *testing.T) {
	small := &m3u8.Variant{URI: "/sample/small.m3u8", VariantParams: m3u8.VariantParams{Bandwidth: 10}}
	medium := &m3u8.Variant{URI: "/sample/medium.m3u8", VariantParams: m3u8.VariantParams{Bandwidth: 20}}
	big := &m3u8.Variant{URI: "/sample/big.m3u8", VariantParams: m3u8.VariantParams{Bandwidth: 30}}

	tests := []struct {
		masterpl *m3u8.MasterPlaylist
	}{
		{&m3u8.MasterPlaylist{Variants: []*m3u8.Variant{small, medium, big}}},
		{&m3u8.MasterPlaylist{Variants: []*m3u8.Variant{small, big, medium}}},
		{&m3u8.MasterPlaylist{Variants: []*m3u8.Variant{big, medium, small}}},
	}

	for i, tt := range tests {
		res := findBiggestVideo(tt.masterpl)
		if res != big {
			t.Errorf("%d: big is not returned.", i)
		}
	}
}

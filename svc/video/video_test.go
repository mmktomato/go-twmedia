package video

import (
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/grafov/m3u8"
	mock_extcmd "github.com/mmktomato/go-twmedia/svc/extcmd/_mock"
	"github.com/mmktomato/go-twmedia/util"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func newVideoServiceImplForTest(t *testing.T, cb func(*VideoServiceImpl, *mock_extcmd.MockExternalCmdService)) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock_extcmd.NewMockExternalCmdService(mockCtrl)
	svc := NewVideoServiceImpl(mock, &util.HttpClient{}, util.NewTinyLogger(false))

	cb(svc, mock)
}

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
		newVideoServiceImplForTest(t, func(svc *VideoServiceImpl, _ *mock_extcmd.MockExternalCmdService) {
			res := svc.findBiggestVideo(tt.variants)
			if res != big {
				t.Errorf("%d: big is not returned.", i)
			}
		})
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
		newVideoServiceImplForTest(t, func(svc *VideoServiceImpl, _ *mock_extcmd.MockExternalCmdService) {
			httpmock.Activate()
			defer httpmock.Deactivate()

			jsurl := "https://localhost/TwitterVideoPlayerIframe.js"
			httpmock.RegisterResponder("GET", jsurl, httpmock.NewStringResponder(200, tt.jsresp))

			token, err := svc.extractAuthToken(jsurl)
			if token != tt.token {
				t.Errorf("%d: unexpected token -> %s", i, token)
				t.Errorf("%d: err -> %v", i, err)
			}
		})
	}
}

func TestFetchTrackInfo(t *testing.T) {
	newVideoServiceImplForTest(t, func(svc *VideoServiceImpl, _ *mock_extcmd.MockExternalCmdService) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		buf, err := ioutil.ReadFile("testdata/trackinfo.json")
		if err != nil {
			t.Fatal(err)
		}

		tweetId := "myTweetId"
		jsonUrl := "https://api.twitter.com/1.1/videos/tweet/config/myTweetId.json"
		resp := string(buf)
		httpmock.RegisterResponder("GET", jsonUrl, httpmock.NewStringResponder(200, resp))

		track, err := svc.fetchTrackInfo(tweetId, "myAuthToken")
		if track.ContentId != "myContentId" {
			t.Errorf("ContentId not match -> %v", track.ContentId)
		}
		if track.PlaylistUrl != "https://example.com/myvideo.m3u8" {
			t.Errorf("PlaylistUrl not match -> %v", track.PlaylistUrl)
		}
	})
}

func TestSavePlaylist(t *testing.T) {
	newVideoServiceImplForTest(t, func(svc *VideoServiceImpl, mockExtCmdService *mock_extcmd.MockExternalCmdService) {
		mockExtCmdService.EXPECT().RunFfmpeg(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		trackInfo := &TrackInfo{"dummyContentId", "http://localhost/master_playlist.m3u8"}

		httpmock.Activate()
		defer httpmock.Deactivate()

		masterPlBuf, err := ioutil.ReadFile("testdata/master_playlist.m3u8")
		if err != nil {
			t.Fatal(err)
		}
		mediaPlBuf, err := ioutil.ReadFile("testdata/media_playlist.m3u8")
		if err != nil {
			t.Fatal(err)
		}

		httpmock.RegisterResponder(
			"GET", "http://localhost/master_playlist.m3u8", httpmock.NewStringResponder(200, string(masterPlBuf)))
		httpmock.RegisterResponder(
			"GET", "http://localhost/media_playlist.m3u8", httpmock.NewStringResponder(200, string(mediaPlBuf)))

		err = svc.SavePlaylist(trackInfo)
		if err != nil {
			t.Error(nil)
		}
	})
}

func TestFindCsrfTokenFromCookie(t *testing.T) {
	tests := []struct {
		cookie, csrfToken string
	}{
		{"foo=FOO; ct0=mytoken; bar=BAR", "mytoken"},
		{"foo=FOO; bar=BAR", ""},
	}

	for i, tt := range tests {
		newVideoServiceImplForTest(t, func(svc *VideoServiceImpl, _ *mock_extcmd.MockExternalCmdService) {
			res := svc.findCsrfTokenFromCookie(tt.cookie)
			if res != tt.csrfToken {
				t.Errorf("%d: csrfToken not match", i)
			}
		})
	}
}

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mmktomato/go-twmedia/svc/extcmd"
	"github.com/mmktomato/go-twmedia/svc/tw"
	"github.com/mmktomato/go-twmedia/svc/video"
	"github.com/mmktomato/go-twmedia/util"
)

var tweetService tw.TweetService = tw.NewTweetServiceImpl()
var videoService video.VideoService = video.NewVideoServiceImpl(
	extcmd.NewExternalCmdServiceImpl())

func onTweetFetched(tweetUrl string, r io.Reader) error {
	tweet, err := tweetService.ParseTweet(tweetUrl, r)
	if err != nil {
		return err
	}

	// TODO: move to svc/image
	for url, filename := range tweet.ImageUrls {
		err := util.Fetch(url, func(r io.Reader) error {
			return util.Save(filename, r)
		})
		if err == nil {
			fmt.Println(filename)
		} else {
			fmt.Printf("%s : %v\n", url, err)
		}
	}

	// TODO: move to svc/video
	if tweet.VideoUrl != "" {
		err := util.Fetch(tweet.VideoUrl, func(r io.Reader) error {
			trackInfo, err := videoService.ParseVideo(tweet.TweetId, r)
			if err != nil {
				return err
			}

			return videoService.SavePlaylist(trackInfo)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no url provided.")
		return
	}

	for i, v := range os.Args {
		if i == 0 {
			continue
		}
		err := util.Fetch(v, func(r io.Reader) error {
			return onTweetFetched(v, r)
		})
		if err != nil {
			fmt.Printf("%s : %v\n", v, err)
		}
	}
}

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mmktomato/go-twmedia/twparser"
	"github.com/mmktomato/go-twmedia/util"
)

func onTweetFetched(tweetUrl string, r io.Reader) error {
	twMedia, err := twparser.ParseTweet(r)
	if err != nil {
		return err
	}

	for url, filename := range twMedia.ImageUrls {
		err := util.Fetch(url, func(r io.Reader) error {
			return util.Save(filename, r)
		})
		if err == nil {
			fmt.Println(filename)
		} else {
			fmt.Printf("%s : %v\n", url, err)
		}
	}

	// TODO: save video
	if twMedia.VideoUrl != "" {
		err := util.Fetch(twMedia.VideoUrl, func(r io.Reader) error {
			trackInfo, err := twparser.ParseVideo(tweetUrl, r)
			if err != nil {
				return err
			}
			fmt.Println(trackInfo)
			return nil
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

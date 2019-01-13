package main

import (
	"io"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mmktomato/go-twmedia/svc/extcmd"
	"github.com/mmktomato/go-twmedia/svc/tw"
	"github.com/mmktomato/go-twmedia/svc/video"
	"github.com/mmktomato/go-twmedia/util"
)

var httpClient util.HttpClient
var logger *util.TinyLogger

var tweetService tw.TweetService
var videoService video.VideoService

type Opts struct {
	Headers    map[string]string `short:"H" long:"header" description:"HTTP header"` // same as curl's one.
	VerboseLog []bool            `short:"v" long:"verbose" description:"verbose log"`
}

func initServices(opts *Opts) {
	httpClient = util.HttpClient{opts.Headers}
	logger = util.NewTinyLogger(0 < len(opts.VerboseLog))
	extcmdService := extcmd.NewExternalCmdServiceImpl(logger)

	videoService = video.NewVideoServiceImpl(extcmdService, &httpClient, logger)
	tweetService = tw.NewTweetServiceImpl(logger)
}

func onTweetFetched(tweetUrl string, r io.Reader) error {
	tweet, err := tweetService.ParseTweet(tweetUrl, r)
	if err != nil {
		return err
	}

	// TODO: move to svc/image
	for url, filename := range tweet.ImageUrls {
		err := httpClient.Fetch(url, func(r io.Reader) error {
			return util.Save(filename, r)
		})
		if err == nil {
			logger.Writeln(filename)
		} else {
			logger.Writef("%s : %v\n", url, err)
		}
	}

	// TODO: move to svc/video
	if tweet.VideoUrl != "" {
		err := httpClient.Fetch(tweet.VideoUrl, func(r io.Reader) error {
			trackInfo, err := videoService.ParseVideo(tweet.TweetId, r)
			if err != nil {
				return err
			}

			logger.Verbosef("track info: %v\n", trackInfo)
			return videoService.SavePlaylist(trackInfo)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getOptions() ([]string, *Opts, error) {
	opts := Opts{}
	parser := flags.NewParser(&opts, flags.IgnoreUnknown)
	args, err := parser.Parse()
	if err != nil {
		return []string{}, nil, err
	}

	for k, v := range opts.Headers {
		opts.Headers[k] = strings.TrimSpace(v)
	}

	return args, &opts, nil
}

func main() {
	args, opts, err := getOptions()
	if err != nil {
		panic(err)
	}

	initServices(opts)

	for _, v := range args {
		if strings.HasPrefix(v, "-") { // unknown option
			continue
		}
		logger.Verbosef("Try: %s\n", v)

		opts.Headers["Referer"] = v

		err := httpClient.Fetch(v, func(r io.Reader) error {
			return onTweetFetched(v, r)
		})
		if err != nil {
			logger.Writef("%s : %v\n", v, err)
		}
	}
}

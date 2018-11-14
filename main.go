package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"strings"

	"github.com/mmktomato/go-twmedia/twparser"
)

func fetch(url string, fn func(io.Reader)) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fn(res.Body)
	return nil
}

func onTweetFetched(r io.Reader) {
	twMedia, err := twparser.ParseTweet(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range twMedia.ImageUrls {
		saveImage(v)
	}

	// TODO: save video
}

func getFilenameFromUrl(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		fmt.Printf("%s : %v\n", url, err)
		return "", err
	}

	tokens := strings.Split(u.Path, "/")
	ret := tokens[len(tokens)-1]
	return strings.Split(ret, ":")[0], nil
}

func saveImage(url string) {
	filename, err := getFilenameFromUrl(url)
	if err != nil {
		fmt.Printf("%s : %v\n", url, err)
		return
	}

	err = fetch(url, func(r io.Reader) {
		buf, readErr := ioutil.ReadAll(r)
		if readErr != nil {
			fmt.Printf("%s : %v\n", url, readErr)
			return
		}

		writeErr := ioutil.WriteFile(filename, buf, 0644)
		if writeErr != nil {
			fmt.Printf("%s : %v\n", url, writeErr)
		}
	})
	if err != nil {
		fmt.Printf("%s : %v\n", url, err)
	}
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
		err := fetch(v, onTweetFetched)
		if err != nil {
			fmt.Printf("%s : %v\n", v, err)
		}
	}
}

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

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

	for url, filename := range twMedia.ImageUrls {
		err = saveImage(url, filename)
		if err == nil {
			fmt.Println(filename)
		} else {
			fmt.Printf("%s : %v\n", url, err)
		}
	}

	// TODO: save video
}

func saveImage(url, filename string) error {
	return fetch(url, func(r io.Reader) {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Printf("%s : %v\n", url, err)
			return
		}

		err = ioutil.WriteFile(filename, buf, 0644)
		if err != nil {
			fmt.Printf("%s : %v\n", url, err)
		}
	})
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

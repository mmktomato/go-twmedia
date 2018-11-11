package main

import (
	"fmt"
	"os"
	"net/http"
	"io"
)

func fetchTweet(url string, fn func(io.Reader)) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fn(res.Body)
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
		fetchErr := fetchTweet(v, func(r io.Reader) {
			twMedia, parseErr := ParseTweet(r)
			if parseErr != nil {
				fmt.Printf("%s : %v\n", v, parseErr)
			}
			fmt.Println(twMedia)
		})
		if fetchErr != nil {
			fmt.Printf("%s : %v\n", v, fetchErr)
			continue
		}
	}
}


package util

import (
	"io"
	"io/ioutil"
	"net/http"
)

func Fetch(url string, fn func(io.Reader) error) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return fn(res.Body)
}

func Save(filename string, r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		return err
	}

	return nil
}

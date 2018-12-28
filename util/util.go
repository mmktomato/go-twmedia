package util

import (
	"io"
	"io/ioutil"
)

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

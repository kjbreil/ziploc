package iniUpdater

import (
	"io/ioutil"
	"os"
)

func File(original, custom, folder string) (*os.File, error) {
	// compare the left and right getting the bytes
	u, err := CompareTwo(original, custom, folder)
	if err != nil {
		return nil, err
	}

	b := u.Bytes()
	tempFile, err := ioutil.TempFile("", "ziploc")
	if err != nil {
		return nil, err
	}
	tempFile.Write(b)
	tempFile.Seek(0, 0)
	return tempFile, nil
}

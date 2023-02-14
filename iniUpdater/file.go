package iniUpdater

import (
	"os"
)

func File(original, custom, folder string) (*os.File, error) {
	// compare the left and right getting the bytes
	u, err := CompareTwo(original, custom, folder)
	if err != nil {
		return nil, err
	}

	b := u.Bytes()
	tempFile, err := os.CreateTemp("", "ziploc")
	if err != nil {
		return nil, err
	}
	tempFile.Write(b)
	tempFile.Seek(0, 0)
	return tempFile, nil
}

func ManualFile(file string, folder string, section string, key string, value string) (*os.File, error) {
	u := CreateManual(file, folder, section, key, value)
	b := u.Bytes()
	tempFile, err := os.CreateTemp("", "ziploc")
	if err != nil {
		return nil, err
	}
	tempFile.Write(b)
	tempFile.Seek(0, 0)
	return tempFile, nil
}

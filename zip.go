package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (o *Option) makeZip() error {
	zipFilename := o.dss.Name + ".zip"

	zipPath := filepath.Join(o.ZipOut, zipFilename)

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	files := make(map[string]os.FileInfo)

	// "walk" the BaseFolder for files
	err = filepath.Walk(o.TempDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files[path] = info
		}
		return nil
	})
	if err != nil {
		return err
	}

	for p, i := range files {
		err = addFileToZip(zipWriter, p, i, o.TempDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string, fileInfo os.FileInfo, tempDir string) error {
	fileData, err := os.Open(filePath)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Name = strings.Replace(filePath, tempDir, "", 1)

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileData)
	if err != nil {
		return err
	}

	return nil
}

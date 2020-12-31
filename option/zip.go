package option

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (o *Option) MakeZip(root string) error {

	name := strings.ReplaceAll(strings.Title(o.Name), " ", "_")

	zipFilename := name + ".zip"

	zipPath := filepath.Join(root, o.ZipOut, zipFilename)

	log.Printf("making zip file %s\n", zipPath)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}

	// open a new zipWriter, explicitly closed at end of function
	zipWriter := zip.NewWriter(zipFile)

	files := make(map[string]os.FileInfo)

	// "walk" the BaseFolder for files
	err = filepath.Walk(filepath.Join(root, o.TempDir), func(path string, info os.FileInfo, err error) error {
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

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	err = zipFile.Close()
	if err != nil {
		return err
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

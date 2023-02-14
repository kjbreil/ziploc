package option

import (
	"fmt"
	"github.com/klauspost/compress/zip"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (o *Option) MakeZip(root string) error {
	caser := cases.Title(language.AmericanEnglish)
	name := strings.ReplaceAll(caser.String(o.Name), " ", "_")

	zipFilename := name + ".zip"

	// check and create ZipOut folder if it does not exist
	zipOutFolder := filepath.Join(root, o.ZipOut)
	if _, err := os.Stat(zipOutFolder); os.IsNotExist(err) {
		err = os.MkdirAll(zipOutFolder, 0777)
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}

	zipPath := filepath.Join(root, o.ZipOut, zipFilename)

	log.Printf("making zip file %s\n", zipPath)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("could not write zip: %s, %v", o.ZipOut, err)
	}

	// open a new zipWriter, explicitly closed at end of function
	zipWriter := zip.NewWriter(zipFile)

	if o.Version != "" {
		_ = zipWriter.SetComment(fmt.Sprintf("V%s %s", o.Version, time.Now().Format("2006-01-02")))
	}

	files := make(map[string]os.FileInfo)

	// "walk" the BaseFolder for files
	err = filepath.Walk(filepath.Join(root, o.tempSubDir()), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files[path] = info
		}
		return nil
	})
	if err != nil {
		return err
	}

	for p, i := range files {
		if i.Name() == "." {
			continue
		}
		err = addFileToZip(zipWriter, p, i, o.tempSubDir())
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
	header.Name = strings.Replace(header.Name, `\`, `/`, -1)
	header.Name = strings.TrimLeft(header.Name, `/`)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileData)
	if err != nil {
		return err
	}

	return fileData.Close()
}

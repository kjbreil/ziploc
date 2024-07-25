package option

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/klauspost/compress/zip"
)

func (o *Option) MakeZip(root string) error {

	name := specialCase(o.Name)
	// name = strings.ReplaceAll(name, "_Ncbp_", "_NCBP_")

	// name := strings.ReplaceAll(o.Name, " ", "_")

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

	zipExists, _ := os.Stat(zipPath)
	if zipExists != nil {
		err := os.Remove(zipPath)
		if err != nil {
			return fmt.Errorf("could not remove zip file: %v", err)
		}
	}
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
		err = addFileToZip(zipWriter, p, i, filepath.Join(o.root, o.tempSubDir()))
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

func specialCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	caser := cases.Title(language.AmericanEnglish)
	s = strings.ReplaceAll(caser.String(s), " ", "_")

	u := []string{
		"NCBP",
		"FMS",
		"GL",
		"TRS",
		"RCV",
		"RSSG",

		"DoorDash",
		"GreatLakes",
		"NoQ",
		"RosieApp",
		"SuperValu",
		"NashFinch",
		"InstaCart",
		"BRdata",
		"iPro",
		"ACNielsen",
	}
	for _, v := range u {
		re := regexp.MustCompile(fmt.Sprintf(`(?i)_%s`, v))
		s = re.ReplaceAllString(s, fmt.Sprintf("_%s", v))
	}
	//
	// c := map[string]string{
	// 	"doordash":   "DoorDash",
	// 	"greatlakes": "GreatLakes",
	// 	"noq":        "NoQ",
	// 	"rosieapp":   "RosieApp",
	// 	"supervalu":  "SuperValu",
	// 	"nashfinch":  "NashFinch",
	// 	"instacart":  "InstaCart",
	// 	"brdata":     "BRdata",
	// 	"ipro":       "iPro",
	// 	"acnielsen":  "ACNielsen",
	// }
	//
	// for k, v := range c {
	//
	// 	re := regexp.MustCompile(fmt.Sprintf(`(?i)_%s`, k))
	// 	s = re.ReplaceAllString(s, fmt.Sprintf("_%s", v))
	// }

	return s
}

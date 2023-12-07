package extract

import (
	"fmt"
	"github.com/kjbreil/ziploc/option"
	"github.com/klauspost/compress/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadZip(zipPath string, basePath string, nestInLast bool, optionMimick *option.Option) (*option.Option, error) {
	// create the option object to make the json file

	_, err := os.Stat(zipPath)
	if err != nil {
		return nil, err
	}
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(r)

	// check if there are files
	if !isValid(r) {
		return nil, fmt.Errorf("nothing in zip file")
	}

	// if mimicking an option build the base dir from the option base dir - 1 dir up
	if optionMimick != nil {
		basePath = filepath.Dir(filepath.Clean(optionMimick.BaseFolder))
	}
	// if nesting in last (usually) add that to the base dir from the last part of the sample
	if nestInLast {
		// log.Println("zp", strings.ReplaceAll(filepath.Base(zipPath), filepath.Ext(zipPath), ""))
		// basePath = filepath.Join(basePath, getSampleNameLast(getSampleName(r)))
		basePath = filepath.Join(basePath, strings.ReplaceAll(filepath.Base(zipPath), filepath.Ext(zipPath), ""))
	}

	o := &option.Option{
		Name:     getSampleName(r),
		Priority: 30,
		IsOption: isOption(r),
		Include: map[string][]string{
			"every": {
				".*",
			},
		},
		Exclude:     nil,
		Ignore:      nil,
		BaseFolder:  basePath,
		TempDir:     "",
		ZipOut:      "",
		Instance:    nil,
		InstanceDir: nil,
		Git:         nil,
		IniMaps:     nil,
	}

	if optionMimick != nil {
		o.ZipOut = optionMimick.ZipOut
		o.TempDir = optionMimick.TempDir
		o.Priority = optionMimick.Priority
	}

	err = extractZip(r, basePath)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func extractZip(r *zip.ReadCloser, dest string) error {
	log.Println("extracting to ", dest)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	for _, f := range r.File {
		// log.Println(getOptionSampleSampleName(r))
		optionSampleName := getOptionSampleSampleName(r)
		fName := f.Name
		outFilePath := filepath.Join(dest, strings.ReplaceAll(fName, optionSampleName, ""))
		if !strings.HasPrefix(outFilePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", outFilePath)
		}
		if f.FileInfo().IsDir() {
			// Make Folder
			err := os.MkdirAll(outFilePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		err = outFile.Close()
		if err != nil {
			return err
		}
		err = rc.Close()
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
	}
	return nil
}

package option

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func (o *Option) Walk(folder string) error {

	stat, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", folder)
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, eachFile := range files {
		p := filepath.Join(folder, eachFile.Name())
		if eachFile.IsDir() {
			err = o.Walk(p)
			if err != nil {
				return err
			}
			continue
		}
		// log.Printf("adding file: %s\n", p)
		o.files[p] = eachFile
	}
	return nil
}

func (o *Option) included(current string) bool {
	b, _ := o.check(current)
	return b
}

func (o *Option) check(current string) (bool, string) {
	if current[0] == []byte(".")[0] {
		return false, ""
	}

	for folder := range o.Include {
		if folder == "EVERY" {
			continue
		}
		for _, item := range o.Include[folder] {
			reg, err := regexp.Compile("(?i)" + item)
			if err != nil {
				log.Panicf("regex could not compile %s: %v", item, err)
			}

			if reg.MatchString(current) {
				return true, folder
			}

		}
	}

	for _, item := range o.Include["every"] {
		reg, err := regexp.Compile("(?i)" + item)
		if err != nil {
			log.Panicf("regex could not compile %s: %v", item, err)
		}

		if reg.MatchString(current) {
			return true, "EVERY"
		}

	}
	return false, ""
}

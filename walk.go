package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func (o *Option) walk(folder string) error {
	if o.excluded(filepath.Base(folder)) {
		return nil
	}

	stat, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", folder)
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, eachFile := range files {
		p := filepath.Join(folder, eachFile.Name())
		if eachFile.IsDir() {
			o.walk(p)
		}
		if !eachFile.IsDir() && o.included(eachFile.Name()) && !o.excluded(p) {
			log.Printf("adding file: %s\n", p)
			o.files[p] = eachFile
		}
	}
	return nil
}

func (o *Option) included(current string) bool {
	b, _ := o.check(current)
	return b
}

func (o *Option) check(current string) (bool, string) {
	for folder := range o.Include {
		for _, item := range o.Include[folder] {
			// check if it matches, uppercase both
			// fmt.Println(current, item, strings.EqualFold(current, item))

			reg, err := regexp.Compile("(?i)" + item)
			if err != nil {
				log.Panicf("regex could not compile %s: %v", item, err)
			}

			if reg.MatchString(current) {
				return true, folder
			}

		}
	}
	return false, ""
}

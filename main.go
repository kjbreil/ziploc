package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjbreil/ziploc/dss"
)

// Option is a LOC option type
type Option struct {
	Name       string              `json:"name"`
	Git        *Git                `json:"git,omitempty"`
	Priority   int                 `json:"priority"`
	BaseFolder string              `json:"base_folder"`
	Include    map[string][]string `json:"include"`
	Exclude    []string            `json:"exclude"`
	Type       string              `json:"type"` // Options or Samples
	TempDir    string              `json:"temp_dir"`
	ZipOut     string              `json:"zip_out"`
	// Not Exported
	files map[string]os.FileInfo
	dss   *dss.DSS
}

type Git struct {
	URL       string `json:"url"`
	Branch    string `json:"branch"`
	Username  string `json:"username"`
	AuthToken string `json:"auth_token"`
}

var (
	configLocation = flag.String("c", "", "Config file")
	makeTemplate   = flag.Bool("template", false, "output a template config file and exit")
)

func main() {

	// Parse the flags
	flag.Parse()

	if *makeTemplate {
		configTemplate(true)
		return
	}

	o, err := readConfig()
	if err != nil {
		log.Panic(err)
	}

	if o.Git != nil {
		err = o.doGitRepo()
		if err != nil {
			log.Panicln(err)
		}
	}

	// TODO: Delete the temp directory before doing anything
	// stat the tempdir and delete if there is no error (directory exists)
	_, err = os.Stat(o.TempDir)
	if err == nil {
		log.Printf("Removing temp directory: %s\n", o.TempDir)
		err = os.RemoveAll(o.TempDir)
		if err != nil {
			log.Panicln(err)
		}
	}

	// TODO: Replace spaces in foldername to make safer and easier to work with

	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	info, err := os.Stat(o.BaseFolder)
	if err != nil || !info.IsDir() {
		log.Panicf("Panicing")
	}

	err = o.walk(o.BaseFolder)

	if err != nil {
		log.Panic(err)
	}

	// Create a new DSS for this run
	o.dss = dss.New(o.Name, o.Priority)
	// loop over the paths and add to the DSS
	for path, info := range o.files {
		// ignore root files
		if o.folder(info.Name()) == "ROOT" {
			continue
		}
		err = o.dss.Add(path)
		if err != nil {
			log.Panicf("error adding %s to dss: %v\n", path, err)
		}
	}

	// Write the DSS to the temp directory
	err = o.dss.Write(filepath.Join(o.TempDir, o.Type, o.dss.Name))
	if err != nil {
		log.Println(err)
	}

	err = o.writeInstall()
	if err != nil {
		log.Panicf("error writing install file: %v\n", err)
	}

	err = o.copyFiles()
	if err != nil {
		log.Panic(err)
	}

	err = o.makeZip()
	if err != nil {
		log.Panic(err)
	}

}

func (o *Option) folder(current string) string {
	_, s := o.check(current)
	return strings.ToUpper(s)
}

func (o *Option) safeName() string {
	return strings.ReplaceAll(strings.Title(o.Name), " ", "_")
}

// makePath uses filepath.Join to safely create the path to the file using OS independent paths
func (o *Option) makePath(folder string, filename string) string {
	p := filepath.Join(o.TempDir, o.Type, o.dss.Name, folder, filename)
	return p
}

func (o *Option) copyFiles() error {
	for ep, ef := range o.files {
		f, err := os.Open(ep)
		if err != nil {
			return err
		}
		var dest string
		folder := o.folder(ef.Name())
		switch folder {
		case "EVERY":
			dest = o.makePath(dss.Destination(ep), ef.Name())
		case "ROOT":
			dest = o.makePath("", ef.Name())
		default:
			newPath := filepath.Join(folder, dss.Destination(ep))
			dest = o.makePath(newPath, ef.Name())
		}

		if _, err := os.Stat(filepath.Dir(dest)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dest), 0777)
			if err != nil {
				return fmt.Errorf("could not create directory: %v", err)
			}
		}

		o, err := os.Create(dest)
		if err != nil {
			return err
		}
		_, err = io.Copy(o, f)
		if err != nil {
			return err
		}
		_ = f.Close()
		_ = o.Close()
	}
	return nil
}

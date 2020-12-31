package main

import (
	"flag"
	"github.com/kjbreil/ziploc/gui"
	"github.com/kjbreil/ziploc/option"
	"log"
	"os"
	"path/filepath"

	"github.com/kjbreil/ziploc/dss"
)

var (
	configLocation = flag.String("c", "", "Config file")
	makeTemplate   = flag.Bool("template", false, "output a template config file and exit")
)

func main() {

	// Parse the flags
	flag.Parse()

	if *makeTemplate {
		option.ConfigTemplate(true)
		return
	}

	// no configuration passed so run GUI
	if *configLocation == "" {
		gui.OpenGui()
		return
	}

	o, err := option.ReadConfig(configLocation)
	if err != nil {
		log.Panic(err)
	}

	if o.Git != nil {
		err = o.DoGitRepo()
		if err != nil {
			log.Panicln(err)
		}
	}

	// TODO: Delete the temp directory before doing anything
	// stat the tmpdir and delete if there is no error (directory exists)
	_, err = os.Stat(o.TempDir)
	if err == nil {
		log.Printf("Removing temp directory: %s\n", o.TempDir)
		err = os.RemoveAll(o.TempDir)
		if err != nil {
			log.Panicln(err)
		}
	}

	// TODO: Replace spaces in folder name to make safer and easier to work with

	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	info, err := os.Stat(o.BaseFolder)
	if err != nil || !info.IsDir() {
		log.Panicf("Panicing")
	}

	err = o.Walk(o.BaseFolder)

	if err != nil {
		log.Panic(err)
	}

	// Create a new DSS for this run
	o.Dss = dss.New(o.Name, o.Priority)
	// loop over the paths and add to the DSS
	for path, info := range o.Files {
		// ignore root files
		if o.Folder(info.Name()) == "ROOT" {
			continue
		}
		err = o.Dss.Add(path)
		if err != nil {
			log.Panicf("error adding %s to dss: %v\n", path, err)
		}
	}

	// Write the DSS to the temp directory
	err = o.Dss.Write(filepath.Join(o.TempDir, o.Type, o.Dss.Name))
	if err != nil {
		log.Println(err)
	}

	err = o.WriteInstall()
	if err != nil {
		log.Panicf("error writing install file: %v\n", err)
	}

	err = o.CopyFiles()
	if err != nil {
		log.Panic(err)
	}

	err = o.MakeZip()
	if err != nil {
		log.Panic(err)
	}

}

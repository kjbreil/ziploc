package main

import (
	"flag"
	"github.com/kjbreil/ziploc/extract"
	"github.com/kjbreil/ziploc/option"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var (
	configLocation = flag.String("c", "", "Config file")
	makeTemplate   = flag.Bool("template", false, "output a template config file and exit")

	fromZip  = flag.String("zip", "", "read from a zip and build config based on said zip, this is the zip file ot extract")
	basePath = flag.String("base", "base", "base path to which extraction of the zip happens, if a config is provided this is ignored")
)

func main() {
	flag.Parse()

	if *makeTemplate {
		option.ConfigTemplate(true)
		return
	}

	// from zip argument passed
	if *fromZip != "" {
		doFromZip()
		return
	}

	// // no configuration passed so run GUI
	if *configLocation == "" {
		// gui.OpenGui()
		log.Println("no config specified")
		return
	}
	cl := filepath.Clean(*configLocation)

	if filepath.Base(cl) == "*" {
		folder, err := os.Stat(filepath.Dir(cl))
		if err != nil {
			panic(err)
		}
		if !folder.IsDir() {
			panic("something went wrong, not a folder")
		}
		log.Println(filepath.Dir(cl))
		filepath.Walk(filepath.Dir(cl), func(path string, info fs.FileInfo, err error) error {

			if filepath.Ext(path) != ".json" {
				return nil
			}
			doSingleConfig(path)

			return nil
		})
	} else {
		doSingleConfig(*configLocation)
	}
}

func doSingleConfig(configPath string) {
	o, err := option.ReadConfig(configPath)
	if err != nil {
		log.Panic(err)
	}

	if o.Git != nil {
		err = o.DoGitRepo()
		if err != nil {
			log.Panicln(err)
		}
	}

	// get the files
	err = o.GetBaseFiles()
	if err != nil {
		log.Panic(err)
	}
	// make dss object for base files
	err = o.GetBaseDSS()
	if err != nil {
		log.Panic(err)
	}
	// if instanceDir is set walk it
	if o.InstanceDir != nil {
		err = o.WalkInstance(*o.InstanceDir)
		if err != nil {
			log.Panic(err)
		}
	}

	// find any INI's and add them to the map
	err = o.FindINI()
	if err != nil {
		log.Panic(err)
	}

	if o.InstanceDir != nil {
		err = o.FromInstance()
		if err != nil {
			log.Panic(err)
		}
	}

	err = o.FromBase("")
	if err != nil {
		log.Panic(err)
	}

	o.WriteConfig(configPath)
}

func doFromZip() {
	fz := filepath.Clean(*fromZip)

	if filepath.Base(fz) == "*" {
		folder, err := os.Stat(filepath.Dir(fz))
		if err != nil {
			panic(err)
		}
		if !folder.IsDir() {
			panic("something went wrong, not a folder")
		}
		log.Println(filepath.Dir(fz))
		filepath.Walk(filepath.Dir(fz), func(path string, info fs.FileInfo, err error) error {

			if filepath.Ext(path) != ".zip" {
				return nil
			}
			doSingleZip(path)

			return nil
		})

	} else {

		doSingleZip(*fromZip)
	}
	return
}

func doSingleZip(path string) {
	log.Println("doing single zip")
	var err error
	var fromOption *option.Option
	var newOption *option.Option

	// config exists so make o from that
	if *configLocation != "" {
		fromOption, err = option.ReadConfig(*configLocation)
		if err != nil {
			log.Panicln(err)
		}
	}
	newOption, err = extract.ReadZip(path, *basePath, true, fromOption)
	if err != nil {
		log.Panicln(err)
	}
	// get the folder that the fromOption is in
	configLocation := filepath.Dir(*configLocation)

	newOption.WriteConfig(filepath.Join(configLocation, newOption.Name) + ".json")
}

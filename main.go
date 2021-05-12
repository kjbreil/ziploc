package main

import (
	"flag"
	"github.com/kjbreil/ziploc/option"
	"log"
)

var (
	configLocation = flag.String("c", "", "Config file")
	makeTemplate   = flag.Bool("template", false, "output a template config file and exit")
)

func main() {

	// gui.OpenGui()
	// return

	// Parse the flags
	flag.Parse()

	if *makeTemplate {
		option.ConfigTemplate(true)
		return
	}

	// // no configuration passed so run GUI
	if *configLocation == "" {
		// gui.OpenGui()
		log.Println("no config specified")
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

	o.WriteConfig(*configLocation)

}

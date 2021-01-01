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

	// Parse the flags
	flag.Parse()

	if *makeTemplate {
		option.ConfigTemplate(true)
		return
	}

	// no configuration passed so run GUI
	if *configLocation == "" {
		// gui.OpenGui()
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

	if o.InstanceDir != nil {
		err = o.FromInstance()
		if err != nil {
			log.Panic(err)
		}
		// return
	}
	err = o.FromBase("")
	if err != nil {
		log.Panic(err)
	}
}

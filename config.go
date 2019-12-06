package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func readConfig() (*Option, error) {

	if *configLocation == "" {
		return nil, fmt.Errorf("config file not set")

	}

	if filepath.Ext(*configLocation) != ".json" {
		return nil, fmt.Errorf("config file is not json")
	}

	b, err := ioutil.ReadFile(*configLocation)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %v", *configLocation, err)
	}

	// make the option object and unmarshal the json into it
	var o Option
	err = json.Unmarshal(b, &o)

	// make the files map
	o.files = make(map[string]os.FileInfo)

	return &o, err
}

func configTemplate() {
	// setup Option information
	o := Option{
		Name:       "SOME SAMPLE",
		BaseFolder: "c:\\storeman\\",
		files:      make(map[string]os.FileInfo),
		Type:       "Samples",
		TempDir:    "SOME_SAMPLE",
	}
	// just for now until config is built
	o.Include = []string{
		"system.ini",
		"samples.ini",
	}

	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile("./template.json", b, 0666)
	if err != nil {
		log.Panic(err)
	}
}

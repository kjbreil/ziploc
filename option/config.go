package option

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ReadConfig(configLocation *string) (*Option, error) {

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
	o.Files = make(map[string]os.FileInfo)

	// check if the priority isn't filled and default to 30
	if o.Priority == 0 {
		o.Priority = 30
	}

	// create the default excludes
	if len(o.Exclude) == 0 {
		o.Exclude = append(o.Exclude, []string{
			"TEMP",
			"SAMPLES",
			"OPTIONS",
		}...)
	}

	return &o, err
}

func ConfigTemplate(withGit bool) {
	// setup Option information
	o := Option{
		Name:       "SOME SAMPLE",
		Priority:   30,
		BaseFolder: "c:\\storeman\\",
		Files:      make(map[string]os.FileInfo),
		Type:       "Samples",
		TempDir:    "SOME_SAMPLE",
		Include:    make(map[string][]string),
	}

	if withGit {
		o.Git = &Git{
			URL:    "https://github.com/kjbreil/ncbp_lab",
			Branch: "master",
		}
	}

	o.Include["every"] = []string{
		"system.ini",
		"samples.ini",
	}
	o.WriteConfig("./template.json")
}

func (o *Option) WriteConfig(path string) {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(path, b, 0666)
	if err != nil {
		log.Panic(err)
	}
}

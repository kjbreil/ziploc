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
	o.files = make(map[string]os.FileInfo)
	o.instanceFiles = make(map[string]os.FileInfo)

	// set the root based on the config
	o.root = filepath.Dir(*configLocation)

	// check if the priority isn't filled and default to 30
	if o.Priority == 0 {
		o.Priority = 30
	}

	return &o, err
}

func (o *Option) DefaultExclude() {
	o.Exclude = &[]string{
		"TEMP",
		"SAMPLES",
		"OPTIONS",
	}
}

func (o *Option) DefaultInclude() {
	o.Include = make(map[string][]string)
	o.Include["every"] = []string{".*"}
}

func ConfigTemplate(withGit bool) {
	// setup Option information
	o := Option{
		Name:       "SOME SAMPLE",
		Priority:   30,
		BaseFolder: "c:\\storeman\\",
		files:      make(map[string]os.FileInfo),
		IsSample:   true,
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

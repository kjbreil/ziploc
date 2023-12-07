package option

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SMSX struct {
	XMLName        xml.Name `xml:"smsx"`
	ProjectName    string
	Creator        string
	SamplePriority int
	IsOption       int `xml:"isOption"`
}

func ReadSmsxConfig(configPath string) (*Option, error) {

	if configPath == "" {
		return nil, fmt.Errorf("config file not set")

	}

	if filepath.Ext(configPath) != ".smsx" {
		return nil, fmt.Errorf("config file is not json")
	}

	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %v", configPath, err)
	}

	// make the option object and unmarshal the json into it
	var smsx SMSX

	smsx.Creator = "NCBP"

	err = xml.Unmarshal(b, &smsx)

	var isOption bool
	if smsx.IsOption == 1 {
		isOption = true
	}

	o := Option{
		Name:           smsx.ProjectName,
		Priority:       smsx.SamplePriority,
		IsOption:       isOption,
		Include:        map[string][]string{"every": []string{".*"}},
		Exclude:        nil,
		Ignore:         nil,
		BaseFolder:     "./SmsCode",
		TempDir:        "Temp",
		ZipOut:         "Release",
		Version:        "",
		Instance:       nil,
		InstanceDir:    nil,
		Git:            nil,
		IniMaps:        nil,
		root:           "",
		files:          nil,
		instanceFiles:  nil,
		dss:            nil,
		configLocation: strings.ReplaceAll(configPath, "smsx", "json"),
	}

	// make the files map
	o.files = make(map[string]os.FileInfo)
	o.instanceFiles = make(map[string]os.FileInfo)

	// set the root based on the config unless its a network drive
	if o.BaseFolder[:2] != "\\\\" {

		o.root = filepath.Dir(configPath)
	}

	// check if the priority isn't filled and default to 30
	if o.Priority == 0 {
		o.Priority = 30
	}

	return &o, err
}

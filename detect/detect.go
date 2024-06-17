package detect

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/iniUpdater"
	"os"
	"path/filepath"
)

type Detect struct {
	ignores     []string
	optionsDirs []string
	samplesDirs []string
	instanceDir string
	includes    []string
	originalDSS *dss.DSS
	customDSS   *dss.DSS
	inis        map[string][]*iniUpdater.U
}

var ValidDirsInOffice map[string]struct{}
var ValidDirsInStoreman map[string]struct{}

func init() {
	ValidDirsInOffice = map[string]struct{}{
		"usp":  {},
		"ssm":  {},
		"sqr":  {},
		"rtm":  {},
		"load": {},
		"lbz":  {},
		"htm":  {},
		"cgi":  {},
	}
	ValidDirsInStoreman = map[string]struct{}{
		"office": {},
	}
}

func New(instanceDir string) *Detect {

	ignores := []string{
		".git",
		".vscode",
		".DS_Store",
	}
	includes := make([]string, 0, len(ValidDirsInOffice)+len(ValidDirsInStoreman))
	for eachDir := range ValidDirsInOffice {
		includes = append(includes, eachDir)
	}
	for eachDir := range ValidDirsInStoreman {
		includes = append(includes, eachDir)
	}
	return &Detect{
		ignores:     ignores,
		includes:    includes,
		instanceDir: instanceDir,
	}
}

func (d *Detect) AddIgnore(ignore ...string) {
	d.ignores = append(d.ignores, ignore...)
}

func (d *Detect) AddOptionDir(optionDir ...string) {
	d.optionsDirs = append(d.optionsDirs, optionDir...)
}

func (d *Detect) AddOptionDirs(optionsDir string) error {
	stat, err := os.Stat(optionsDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", optionsDir)
	}

	files, err := os.ReadDir(optionsDir)
	if err != nil {
		return err
	}
	for _, eachFile := range files {
		p := filepath.Join(optionsDir, eachFile.Name())
		if eachFile.IsDir() {
			d.optionsDirs = append(d.optionsDirs, p)
		}

	}
	return nil
}

func (d *Detect) AddSampleDirs(samplesDir string) error {
	stat, err := os.Stat(samplesDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", samplesDir)
	}

	files, err := os.ReadDir(samplesDir)
	if err != nil {
		return err
	}
	for _, eachFile := range files {
		p := filepath.Join(samplesDir, eachFile.Name())
		if eachFile.IsDir() {
			d.samplesDirs = append(d.samplesDirs, p)
		}

	}
	return nil
}

func (d *Detect) AddSampleDir(sampleDir ...string) {
	d.samplesDirs = append(d.samplesDirs, sampleDir...)
}

func (d *Detect) Detect() error {
	if d.originalDSS == nil {
		d.originalDSS = dss.New("original", 30)
		var err error
		for _, optionDir := range d.optionsDirs {
			nd := dss.New(filepath.Base(optionDir), 30)
			nd.Author = optionDir
			err = nd.WalkDirProgress(optionDir)
			if err != nil {
				return err
			}
			d.originalDSS.Merge(nd)
		}
		for _, sampleDir := range d.samplesDirs {
			nd := dss.New(filepath.Base(sampleDir), 25)
			nd.Author = sampleDir
			err = nd.WalkDirProgress(sampleDir)
			if err != nil {
				return err
			}
			d.originalDSS.Merge(nd)
		}
	}
	if d.customDSS == nil {
		d.customDSS = dss.New("custom", 1)
		d.customDSS.Author = d.instanceDir
		d.customDSS.Ignore = d.ignores
		d.customDSS.Includes = d.includes
		err := d.customDSS.WalkDirProgress(d.instanceDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Detect) Compare() (dss.Entries, error) {
	err := d.Detect()
	if err != nil {
		return nil, err
	}

	var entries dss.Entries
	for _, e := range d.customDSS.Data {
		if !d.originalDSS.Matches(e) {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

func (d *Detect) INIs() map[string][]*iniUpdater.U {
	return d.inis
}

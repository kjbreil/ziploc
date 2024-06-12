package detect

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"os"
	"path/filepath"
)

type Detect struct {
	ignores     []string
	optionsDirs []string
	samplesDirs []string
	instanceDir string
}

func New(instanceDir string) *Detect {
	return &Detect{
		ignores: []string{
			".git",
			".vscode",
			".DS_Store",
		},
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

func (d *Detect) Compare() (dss.Entries, error) {
	originalDss := dss.New("original", 30)
	var err error
	for _, optionDir := range d.optionsDirs {
		nd := dss.New(filepath.Base(optionDir), 30)

		err = nd.WalkDirProgress(optionDir)
		if err != nil {
			return nil, err
		}
		originalDss.Merge(nd)
	}
	for _, sampleDir := range d.samplesDirs {
		nd := dss.New(filepath.Base(sampleDir), 25)
		err = nd.WalkDirProgress(sampleDir)
		if err != nil {
			return nil, err
		}
		originalDss.Merge(nd)
	}

	customDss := dss.New("custom", 1)
	customDss.Ignore = d.ignores
	err = customDss.WalkDirProgress(d.instanceDir)

	var entries dss.Entries
	for _, e := range customDss.Data {
		if !originalDss.Matches(e) {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

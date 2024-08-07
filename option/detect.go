package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/detect"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/macro"
	"os"
	"path/filepath"
	"strings"
)

type Detect struct {
	StoremanDir  string   `json:"storeman_dir,omitempty"`
	ChangedFiles []string `json:"changed_files,omitempty"`
	Ignores      []string `json:"ignores,omitempty"`
}

func (o *Option) DetectChangesFromInstance(instanceDir string, baseDir string) error {
	var err error
	d := detect.New(instanceDir)
	d.ResetIncludes()
	d.AddIgnore(".*FirstLoad.*", ".*XchDev.*")
	if len(o.Include) > 0 {
		for _, s := range o.Include {
			d.AddInclude(s...)
		}

	}

	err = d.AddOptionDirs(filepath.Join(baseDir))
	if err != nil {
		return err
	}

	entries, err := d.Compare()
	if err != nil {
		return err
	}

	changedFiles := make([]string, 0, len(entries))

	for _, eachEntry := range entries {
		changedFiles = append(changedFiles, eachEntry.Script)
	}

	o.Include["every"] = changedFiles

	return nil
}

func (o *Option) LookupCustomFiles() error {
	if o.Detect == nil {
		return nil
	}
	var err error
	d := detect.New(filepath.Join(o.Detect.StoremanDir, "Office"))

	err = d.AddOptionDirs(filepath.Join(o.Detect.StoremanDir, "Options"))
	if err != nil {
		return err
	}
	err = d.AddSampleDirs(filepath.Join(o.Detect.StoremanDir, "Samples"))
	if err != nil {
		return err
	}
	err = d.GetINIs()

	entries, err := d.Compare()
	if err != nil {
		return err
	}
	// for _, e := range entries {
	// 	if strings.EqualFold(filepath.Ext(e.Script), ".INI") {
	// 		fmt.Println("here")
	// 	}
	// }

	o.Detect.ChangedFiles = make([]string, 0, len(entries))

	for _, eachEntry := range entries {
		o.Detect.ChangedFiles = append(o.Detect.ChangedFiles, eachEntry.Script)
	}

	o.CopyCustomFiles(entries)
	for k, i := range d.INIs() {
		noExt := strings.ToUpper(k[:len(k)-len(filepath.Ext(k))])
		if len(i) > 0 {

			dest := o.makeBuildPath(o.root, "INBOX", noExt+"_INI.SQL")
			if _, err = os.Stat(filepath.Dir(dest)); os.IsNotExist(err) {
				err = os.MkdirAll(filepath.Dir(dest), 0777)
				if err != nil {
					return fmt.Errorf("could not create directory: %v", err)
				}
			}
			os.WriteFile(dest, i[0].Bytes(), 0666)
		}
	}

	return nil
}

func (o *Option) CopyCustomFiles(entries dss.Entries) error {
	if o.Detect == nil {
		return nil
	}
	for _, eachFile := range entries {
		p := detect.BuildFilePath(o.Detect.StoremanDir, eachFile.Destination, eachFile.Script, false)
		_, err := os.Stat(p)
		if err != nil {
			o.Log().Error(fmt.Sprintf("could not stat %s", p), "err", err)
			continue
		}
		dest := o.makeBuildPath(o.root, eachFile.Destination, eachFile.Script)
		if _, err = os.Stat(filepath.Dir(dest)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dest), 0777)
			if err != nil {
				return fmt.Errorf("could not create directory: %v", err)
			}
		}
		out, err := os.Create(dest)
		if err != nil {
			if out != nil {
				_ = out.Close()
			}
			return err
		}

		f, err := os.Open(p)
		err = macro.Correct(out, f, false)
		_ = out.Close()
		_ = f.Close()
		if err != nil {
			return err
		}

	}
	return nil
}

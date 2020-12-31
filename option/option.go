package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/macro"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Option is a LOC option type
type Option struct {
	Name        string              `json:"name"`
	Git         *Git                `json:"git,omitempty"`
	Priority    int                 `json:"priority"`
	BaseFolder  string              `json:"base_folder"`
	Include     map[string][]string `json:"include"`
	Exclude     []string            `json:"exclude"`
	Type        string              `json:"type"` // Options or Samples
	IsSample    bool                `json:"is_sample"`
	Instance    string              `json:"instance"`
	InstanceDir string              `json:"instance_dir"`
	TempDir     string              `json:"temp_dir"`
	ZipOut      string              `json:"zip_out"`
	// Not Exported
	Files map[string]os.FileInfo
	Dss   *dss.DSS
}

type Git struct {
	URL       string `json:"url"`
	Branch    string `json:"branch"`
	Username  string `json:"username"`
	AuthToken string `json:"auth_token"`
}

func (o *Option) Folder(current string) string {
	_, s := o.check(current)
	return strings.ToUpper(s)
}

func (o *Option) safeName() string {
	return strings.ReplaceAll(strings.Title(o.Name), " ", "_")
}

// makePath uses filepath.Join to safely create the path to the file using OS independent paths
func (o *Option) makePath(folder string, filename string) string {
	p := filepath.Join(o.TempDir, o.Type, o.Dss.Name, folder, filename)
	return p
}

func (o *Option) CopyFiles() error {
	for ep, ef := range o.Files {
		f, err := os.Open(ep)
		if err != nil {
			return err
		}
		var dest string

		m := regexp.MustCompile(`SP-\d{3}`)
		folder := m.FindString(ep)
		if folder == "" {
			folder = o.Folder(ef.Name())
		}

		switch folder {
		case "EVERY":
			dest = o.makePath(dss.Destination(ep), ef.Name())
		case "ROOT":
			dest = o.makePath("", ef.Name())
		default:
			newPath := filepath.Join(folder, dss.Destination(ep))
			dest = o.makePath(newPath, ef.Name())
		}

		if _, err := os.Stat(filepath.Dir(dest)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(dest), 0777)
			if err != nil {
				return fmt.Errorf("could not create directory: %v", err)
			}
		}

		out, err := os.Create(dest)
		if err != nil {
			return err
		}

		err = macro.Correct(out, f)
		if err != nil {
			log.Println(dest)
			return err
		}

		_ = f.Close()
		_ = out.Close()
	}
	return nil
}

func (o *Option) DeleteTemp() {
	// stat the tmpdir and delete if there is no error (directory exists)
	_, err := os.Stat(o.TempDir)
	if err == nil {
		log.Printf("Removing temp directory: %s\n", o.TempDir)
		err = os.RemoveAll(o.TempDir)
		if err != nil {
			log.Panicln(err)
		}
	}
}

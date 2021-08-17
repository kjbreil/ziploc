package option

import (
	"github.com/kjbreil/ziploc/dss"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Option is a LOC option type
type Option struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	IsOption bool   `json:"is_option"`

	Include map[string][]string `json:"include,omitempty"`
	Exclude []string            `json:"exclude,omitempty"`
	Ignore  []string            `json:"ignore,omitempty"`

	BaseFolder string `json:"base_folder"`
	TempDir    string `json:"temp_dir"`
	ZipOut     string `json:"zip_out"`

	Instance    *string `json:"instance,omitempty"`
	InstanceDir *string `json:"instance_dir,omitempty"`

	Git *Git `json:"git,omitempty"`

	IniMaps map[string]IniMap `json:"ini_maps,omitempty"`

	// Not Exported
	root          string
	files         map[string]os.FileInfo
	instanceFiles map[string]os.FileInfo
	dss           *dss.DSS
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

func (o *Option) getType() string {
	t := "Options"
	if !o.IsOption {
		t = "Samples"
	}
	return t
}

// makePath uses filepath.Join to safely create the path to the file using OS independent paths
func (o *Option) makePath(root string, folder string, filename string) string {
	p := filepath.Join(root, o.tempSubDir(), o.getType(), o.dss.Name, folder, filename)
	return p
}

func (o *Option) makeBuildPath(root string, folder string, filename string) string {
	p := filepath.Join(root, o.BaseFolder, folder, filename)
	return p
}

func (o *Option) DeleteTemp(root string) {
	// stat the tmpdir and delete if there is no error (directory exists)
	_, err := os.Stat(filepath.Join(root, o.tempSubDir()))
	if err == nil {
		log.Printf("Removing temp directory: %s\n", filepath.Join(root, o.tempSubDir()))
		err = os.RemoveAll(filepath.Join(root, o.tempSubDir()))
		if err != nil {
			log.Panicln(err)
		}
	}
}

func (o *Option) tempSubDir() string {
	return filepath.Join(o.TempDir, o.Name)
}

package option

import (
	"github.com/kjbreil/ziploc/dss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"log/slog"
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
	Version    string `json:"version"`

	Instance    *string `json:"instance,omitempty"`
	InstanceDir *string `json:"instance_dir,omitempty"`

	SignedDSS *string `json:"signed_dss,omitempty"`

	Git *Git `json:"git,omitempty"`

	IniMaps map[string]IniMap `json:"ini_maps,omitempty"`

	Detect *Detect `json:"detect,omitempty"`

	// Not Exported
	root           string
	files          map[string]os.FileInfo
	instanceFiles  map[string]os.FileInfo
	dss            *dss.DSS
	configLocation string
	logger         *slog.Logger
}

type Git struct {
	URL       string `json:"url"`
	Branch    string `json:"branch"`
	Username  string `json:"username"`
	AuthToken string `json:"auth_token"`
}

func (o *Option) Log() *slog.Logger {
	if o.logger == nil {
		o.logger = slog.Default()
	}
	return o.logger
}

func (o *Option) Folder(current string) string {
	_, s := checkIncluded(current, o.Include)
	return strings.ToUpper(s)
}

func (o *Option) safeName() string {
	caser := cases.Title(language.AmericanEnglish)
	return strings.ReplaceAll(caser.String(o.Name), " ", "_")
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
	p := filepath.Join(root, o.tempSubDir(), o.getType(), o.Name, folder, filename)
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

func (o *Option) SetConfigLocation(name string) {
	o.configLocation = name
}

// Clean updates folders to be clean as per filepant.Clean()
// removes .. and adjusts for os version
func (o *Option) Clean() {
	o.BaseFolder = filepath.Clean(o.BaseFolder)
	o.TempDir = filepath.Clean(o.TempDir)
	o.ZipOut = filepath.Clean(o.ZipOut)
	o.Name = specialCase(o.Name)
}

func (o *Option) ChangeRoot(fp string) {
	o.root = fp
}

func (o *Option) SingleDSS(path string) (*dss.DSS, error) {
	d := dss.New(o.Name, o.Priority)

	err := d.Add(path)
	if err != nil {
		return nil, err
	}

	return d, nil
}

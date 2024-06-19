package dss

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/kjbreil/crcloc"

	"github.com/kjbreil/sil"
)

// dss is fed the info about the files and outputs a dss file
type DSS struct {
	name     string
	SIL      *sil.SIL
	Data     Entries
	Ignore   []string
	priority int
	Includes []string

	m      *sync.Mutex
	Author string
}

func New(name string, priority int) *DSS {
	// make the sil object
	s := sil.Make("dss", Table{})
	// change to a LOAD file
	s.View.Action = "LOAD"
	// append the
	activate(s)

	return &DSS{
		name:     makeName(name),
		SIL:      s,
		priority: priority,
		Author:   "NCBP",
		m:        &sync.Mutex{},
	}
}

func makeName(name string) string {
	return strings.ToUpper(strings.ReplaceAll(name, " ", "_"))
}

// add a file to the dss, fp is the path of the file
func (d *DSS) add(fp string) error {

	// ignore anything in the INBOX as that is not confirmed in the DSS
	if Destination(fp) == "INBOX" {
		return nil
	}

	// ignore anything specified in the ignore list
	for _, ignore := range d.Ignore {
		reg, err := regexp.Compile("(?i)" + ignore)
		if err != nil {
			// TODO: report error
			continue
		}

		if reg.MatchString(fp) {
			return nil
		}
	}
	if len(d.Includes) > 0 {
		var match bool
		for _, include := range d.Includes {
			if strings.ToLower(filepath.Base(filepath.Dir(fp))) == include {
				match = true
			}
		}
		if !match {
			return nil
		}
	}

	f, err := os.Open(fp)
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	if err != nil {
		return fmt.Errorf("could not open %s with error: %v", fp, err)
	}

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("could not stat %s with error: %v", fp, err)
	}

	t := Table{
		Priority:    d.priority,
		Author:      d.Author,
		Option:      d.name,
		Destination: Destination(fp),
		Script:      filepath.Base(fp),
		FileDate:    sil.JulianDateTime(info.ModTime()),
		Signature:   crcloc.HashReader(f),
	}
	d.m.Lock()
	d.Data = append(d.Data, &t)
	d.SIL.View.Data = append(d.SIL.View.Data, t)
	d.m.Unlock()

	return nil
}

// Write writes a dss to file
func (d *DSS) Write(folder string) error {

	fullpath := path.Join(folder, "Inbox", fmt.Sprintf("DSS_%s.sql", d.name))

	if _, err := os.Stat(filepath.Dir(fullpath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(fullpath), 0777)
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}

	err := d.SIL.Write(fullpath, true, false)
	if err != nil {
		return fmt.Errorf("could not write file %s: %v", fullpath, err)
	}

	return nil
}

func activate(s *sil.SIL) {
	s.Append("DELETE FROM DSS_TAB WHERE F2729 IN (SELECT F2729 FROM DSS_LOAD);")
	s.Append("@UPDATE_BATCH(JOB=ADDRPL,TAR=DSS_TAB,KEY=F2729=:F2729 AND F2730=:F2730 AND F2731=:F2731,")
	s.Append("SRC=SELECT * FROM DSS_LOAD);")
	s.Append("DROP TABLE DSS_LOAD;")
}

// Destination takes a filepath and gets the last dir
// exported because the same name is used to make the zip file
func Destination(fp string) string {
	// return the directory where the file is then just the base of the directory, uppercase it
	return strings.ToUpper(filepath.Base(filepath.Dir(fp)))
}

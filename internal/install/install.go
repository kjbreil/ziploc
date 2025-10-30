package install

import (
	"path/filepath"
	"strings"
)

// install file(s) into an instance from a base folder

var inStoreman map[string]struct{}

func init() {
	inStoreman = map[string]struct{}{
		"OFFICE":  {},
		"XCHDEV":  {},
		"XCHFILE": {},
	}
}

func InstallLocation(filePath, instanceFolder string) (string, error) {
	dir := filepath.Base(filepath.Dir(filePath))

	if _, ok := inStoreman[strings.ToUpper(dir)]; ok {
		return filepath.Join(instanceFolder, dir, filepath.Base(filePath)), nil
	}

	return filepath.Join(instanceFolder, "Office", dir, filepath.Base(filePath)), nil
}

func IsTopLevel(filePath, baseDir string) bool {
	return strings.EqualFold(filepath.Dir(filePath), filepath.Clean(baseDir))
}

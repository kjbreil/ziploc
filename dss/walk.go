package dss

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
	"time"
)

func (d *DSS) WalkDir(folder string) error {
	return d.walkDir(folder, nil)
}
func (d *DSS) WalkDirProgress(folder string) error {
	count, err := d.DirFilesCount(folder)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions64(
		int64(count),
		progressbar.OptionSetDescription(folder),
		progressbar.OptionSetWriter(os.Stderr),
		// progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(150*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionEnableColorCodes(true),
	)
	return d.walkDir(folder, bar)
}

func (d *DSS) walkDir(folder string, bar *progressbar.ProgressBar) error {
	stat, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", folder)
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, eachFile := range files {
		p := filepath.Join(folder, eachFile.Name())
		if eachFile.IsDir() {
			err = d.walkDir(p, bar)
			if err != nil {
				return err
			}
			continue
		}
		err = d.add(p)
		if bar != nil {
			bar.Add(1)
		}
		if err != nil {
			return fmt.Errorf("error adding %s to dss: %v\n", p, err)
		}
	}
	return nil
}

func (d *DSS) DirFilesCount(folder string) (int, error) {
	var count int
	stat, err := os.Stat(folder)
	if err != nil {
		return count, err
	}
	if !stat.IsDir() {
		return count, fmt.Errorf("%s is not a directory", folder)
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return count, err

	}
	var newCount int

	for _, eachFile := range files {
		p := filepath.Join(folder, eachFile.Name())
		if eachFile.IsDir() {
			newCount, err = d.DirFilesCount(p)
			if err != nil {
				return count, err
			}
			count += newCount
			continue
		}
		count++

	}
	return count, nil
}

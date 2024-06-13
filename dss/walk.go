package dss

import (
	"context"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"time"
)

func (d *DSS) WalkDir(folder string) error {
	return d.walkDir(folder, nil)
}
func (d *DSS) WalkDirProgress(folder string) error {
	bar := progressbar.NewOptions64(
		int64(100),
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

	paths, err := getPaths(folder)
	if err != nil {
		return err
	}
	if bar != nil {
		bar.ChangeMax(len(paths))
	}

	// addChan := make(chan string, numOfRoutines)

	g, _ := errgroup.WithContext(context.Background())
	numOfRoutines := 100
	g.SetLimit(numOfRoutines)
	for _, p := range paths {
		pp := p
		g.Go(func() error {
			err = d.add(pp)
			if err != nil {
				return err
			}
			if bar != nil {
				err = bar.Add(1)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}

func getPaths(folder string) ([]string, error) {
	stat, err := os.Stat(folder)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", folder)
	}

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var paths []string

	for _, eachFile := range files {
		p := filepath.Join(folder, eachFile.Name())
		if eachFile.IsDir() {
			addPaths, err := getPaths(p)
			if err != nil {
				return nil, err
			}
			paths = append(paths, addPaths...)
			continue
		}
		paths = append(paths, p)
	}
	return paths, nil
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

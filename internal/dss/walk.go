package dss

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
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
	numOfRoutines := 100
	addChan := make(chan string, numOfRoutines)
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	for range numOfRoutines {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for {
				select {
				case p := <-addChan:
					err = d.Add(p)
					if err != nil {
						log.Println(err)
					}
					if bar != nil {
						err = bar.Add(1)
						if err != nil {
							log.Println(err)
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	for _, p := range paths {
		addChan <- p
	}
	for len(addChan) > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	cancel()
	wg.Wait()
	// g, _ := errgroup.WithContext(context.Background())
	// g.SetLimit(numOfRoutines)
	// for _, p := range paths {
	// 	pp := p
	// 	g.Go(func() error {
	// 		err = d.Add(pp)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if bar != nil {
	// 			err = bar.a(1)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 		return nil
	// 	})
	// }
	// if err = g.Wait(); err != nil {
	// 	return err
	// }
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

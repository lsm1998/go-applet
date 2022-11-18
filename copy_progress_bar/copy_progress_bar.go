package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultBufSize = 32 << 10
)

type fileCopyProgressBar struct {
	fileSize int64
	Bar
}

func (c *fileCopyProgressBar) Copy(src, dist string) (err error) {
	if err = filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		c.fileSize += info.Size()
		return nil
	}); err != nil {
		return err
	}
	c.Bar = NewProgressBar(float64(c.fileSize), WithTail(">"), WithFiller("="), WithInterval(time.Second/100))
	if err = c.copy(src, dist); err != nil {
		c.Bar.Cancel()
	}
	return nil
}

func (c *fileCopyProgressBar) copy(src, dist string) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		if err := os.MkdirAll(dist, os.ModePerm); err != nil {
			return err
		}
		fileInfos, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}
		for i := 0; i < len(fileInfos); i++ {
			_src := fmt.Sprintf("%s%c%s", src, os.PathSeparator, fileInfos[i].Name())
			_dist := fmt.Sprintf("%s%c%s", dist, os.PathSeparator, fileInfos[i].Name())
			if err = c.copy(_src, _dist); err != nil {
				return err
			}
		}
		return err
	} else {
		return c.copyFile(src, dist)
	}
}

func (c *fileCopyProgressBar) copyFile(src, dist string) error {
	write, err := os.Create(dist)
	if err != nil {
		return err
	}
	read, err := os.Open(src)
	var buf = make([]byte, defaultBufSize)
	var n int
	for {
		n, err = read.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if _, err = write.Write(buf[:n]); err != nil {
			return err
		}
		c.Bar.Done(float64(n))
	}
	return nil
}

func NewFileCopyProgressBar(opts ...BarOption) Copy {
	return &fileCopyProgressBar{}
}

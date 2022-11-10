package main

import (
	"io"
	"os"
)

type copyProgressBar struct {
}

func (c *copyProgressBar) Copy(src, dist string) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return c.copyDir(src, dist)
	} else {
		return c.copyFile(src, dist)
	}
}

func (c *copyProgressBar) copyDir(src, dist string) error {
	return nil
}

func (c *copyProgressBar) copyFile(src, dist string) error {
	write, err := os.Create(dist)
	if err != nil {
		return err
	}
	read, err := os.Open(src)
	if _, err = io.Copy(write, read); err != nil {
		return err
	}
	return nil
}

func NewCopyProgressBar() Copy {
	// watch, _ := fsnotify.NewWatcher()
	return &copyProgressBar{
		// Watcher: watch,
	}
}

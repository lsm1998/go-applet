package file_sharding

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type fileAdapter struct {
	number   uint
	filename string
	suffix   string
	bufSize  int
}

type AdapterOpt = func(f *fileAdapter)

func (f *fileAdapter) Merge() (string, error) {
	target, err := os.Create(f.filename)
	if err != nil {
		return "", err
	}
	defer target.Close()
	var files = make([]*os.File, 0, f.number)
	for i := 0; i < int(f.number); i++ {
		file, err := os.Open(fmt.Sprintf("%s%s%d", f.filename, f.suffix, i))
		if err != nil {
			return "", err
		}
		defer file.Close()
		files = append(files, file)
	}
	var buf = make([]byte, f.bufSize)
	var n int
	var index int
	var finishCount int
	for {
		if n, err = files[index%len(files)].Read(buf); err == nil {
			if _, err = target.Write(buf[:n]); err != nil {
				return "", err
			}
		} else if errors.Is(io.EOF, err) {
			finishCount++
		}
		if finishCount == len(files) {
			break
		}
		index++
	}
	return f.filename, nil
}

func (f *fileAdapter) Sharding() ([]string, error) {
	src, err := os.Open(f.filename)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	var result = make([]string, 0, f.number)
	var files = make([]*os.File, 0, f.number)
	for i := 0; i < int(f.number); i++ {
		filename := fmt.Sprintf("%s%s%d", f.filename, f.suffix, i)
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		result = append(result, filename)
		files = append(files, file)
	}
	var buf = make([]byte, f.bufSize)
	var n int
	var index int

	for {
		if n, err = src.Read(buf); err == nil {
			_, err = files[index%len(files)].Write(buf[:n])
			index++
		} else if errors.Is(io.EOF, err) {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func WithFilename(filename string) AdapterOpt {
	return func(f *fileAdapter) {
		f.filename = filename
	}
}

func WithNumber(number uint) AdapterOpt {
	return func(f *fileAdapter) {
		f.number = number
	}
}

func NewFileAdapter(opts ...AdapterOpt) Adapter {
	adapter := &fileAdapter{}
	opts = append([]AdapterOpt{
		func(f *fileAdapter) {
			f.number = defaultNumber
			f.suffix = defaultSuffix
			f.bufSize = defaultBufSize
		},
	}, opts...)
	for _, v := range opts {
		v(adapter)
	}
	return adapter
}

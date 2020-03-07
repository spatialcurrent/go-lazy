// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package lazy

import (
	"os"
)

type LazyFile struct {
	name string
	flag int
	perm os.FileMode
	file *os.File
	err  error
}

func NewLazyFile(name string, flag int, perm os.FileMode) *LazyFile {
	return &LazyFile{
		name: name,
		flag: flag,
		perm: perm,
		file: nil,
		err:  nil,
	}
}

func (f *LazyFile) Fd() (uintptr, error) {
	if f.err != nil {
		return 0, f.err
	}
	if f.file == nil {
		file, err := os.OpenFile(f.name, f.flag, f.perm)
		if err != nil {
			f.err = err
			return 0, err
		}
		f.file = file
	}
	return f.file.Fd(), nil
}

func (f *LazyFile) Read(b []byte) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	if f.file == nil {
		file, err := os.OpenFile(f.name, f.flag, f.perm)
		if err != nil {
			f.err = err
			return 0, err
		}
		f.file = file
	}
	return f.file.Read(b)
}

func (f *LazyFile) Stat() (os.FileInfo, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.file == nil {
		file, err := os.OpenFile(f.name, f.flag, f.perm)
		if err != nil {
			f.err = err
			return nil, err
		}
		f.file = file
	}
	fileInfo, err := f.file.Stat()
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func (f *LazyFile) Close() error {
	if f.file == nil {
		return nil
	}
	return f.file.Close()
}

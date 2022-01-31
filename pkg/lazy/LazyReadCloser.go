// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package lazy

import (
	"io"
)

type LazyReadCloser struct {
	init       func() (io.ReadCloser, error)
	readCloser io.ReadCloser
	err        error
}

func (r *LazyReadCloser) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.readCloser == nil {
		readCloser, err := r.init()
		if err != nil {
			r.err = err
			return 0, err
		}
		r.readCloser = readCloser
	}
	return r.readCloser.Read(p)
}

func (r *LazyReadCloser) Close() error {
	if r.readCloser != nil {
		defer func() { r.readCloser = nil }()
		return r.readCloser.Close()
	}
	return nil
}

func NewLazyReadCloser(init func() (io.ReadCloser, error)) *LazyReadCloser {
	return &LazyReadCloser{
		init:       init,
		readCloser: nil,
		err:        nil,
	}
}

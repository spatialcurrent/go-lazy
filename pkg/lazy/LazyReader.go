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

type LazyReader struct {
	init   func() (io.Reader, error)
	reader io.Reader
	err    error
}

func (r *LazyReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.reader == nil {
		reader, err := r.init()
		if err != nil {
			r.err = err
			return 0, err
		}
		r.reader = reader
	}
	return r.reader.Read(p)
}

func (r *LazyReader) Flush() error {
	if flusher, ok := r.reader.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *LazyReader) Close() error {
	if closer, ok := r.reader.(interface{ Close() error }); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func NewLazyReader(init func() (io.Reader, error)) *LazyReader {
	return &LazyReader{
		init:   init,
		reader: nil,
		err:    nil,
	}
}

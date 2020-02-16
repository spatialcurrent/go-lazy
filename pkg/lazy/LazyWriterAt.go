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

type LazyWriterAt struct {
	init   func() (io.WriterAt, error)
	writer io.WriterAt
	err    error
}

func (r *LazyWriterAt) WriteAt(p []byte, off int64) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.writer == nil {
		writer, err := r.init()
		if err != nil {
			r.err = err
			return 0, err
		}
		r.writer = writer
	}
	return r.writer.WriteAt(p, off)
}

func (r *LazyWriterAt) Flush() error {
	if flusher, ok := r.writer.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *LazyWriterAt) Close() error {
	if closer, ok := r.writer.(interface{ Close() error }); ok {
		err := closer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func NewLazyWriterAt(init func() (io.WriterAt, error)) *LazyWriterAt {
	return &LazyWriterAt{
		init:   init,
		writer: nil,
		err:    nil,
	}
}

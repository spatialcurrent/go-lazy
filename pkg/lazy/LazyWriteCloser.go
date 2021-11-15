// =================================================================
//
// Copyright (C) 2021 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package lazy

import (
	"io"
)

type LazyWriteCloser struct {
	init   func() (io.WriteCloser, error)
	writer io.WriteCloser
	err    error
}

func (l *LazyWriteCloser) Write(p []byte) (int, error) {
	if l.err != nil {
		return 0, l.err
	}
	if l.writer == nil {
		writer, err := l.init()
		if err != nil {
			l.err = err
			return 0, err
		}
		l.writer = writer
	}
	return l.writer.Write(p)
}

func (l *LazyWriteCloser) Flush() error {
	if flusher, ok := l.writer.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LazyWriteCloser) Close() error {
	if l.writer == nil {
		return nil
	}
	return l.writer.Close()
}

func NewLazyWriteCloser(init func() (io.WriteCloser, error)) *LazyWriteCloser {
	return &LazyWriteCloser{
		init:   init,
		writer: nil,
		err:    nil,
	}
}

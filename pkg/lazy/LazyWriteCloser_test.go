// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package lazy

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bufferCloser struct {
	*bytes.Buffer
}

func (bufferCloser) Close() error {
	return nil
}

func TestLazyWriteCloser(t *testing.T) {
	buf := bufferCloser{bytes.NewBuffer([]byte{})}
	w := NewLazyWriteCloser(func() (io.WriteCloser, error) {
		return buf, nil
	})
	require.Equal(t, 0, len(buf.Bytes()))
	n, err := w.Write([]byte("hello world"))
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "hello world", buf.String())
	err = w.Close()
	assert.NoError(t, err)
}

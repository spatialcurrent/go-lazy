// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package lazy

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go/aws"
)

func TestLazyWriterAt(t *testing.T) {
	buf := aws.NewWriteAtBuffer(make([]byte, 0))
	w := NewLazyWriterAt(func() (io.WriterAt, error) {
		return buf, nil
	})
	require.Equal(t, 0, len(buf.Bytes()))
	n, err := w.WriteAt([]byte("hello world"), 0)
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "hello world", string(buf.Bytes()))
}

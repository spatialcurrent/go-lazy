// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package lazy

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLazyReader(t *testing.T) {
	in := "hello world"
	r := NewLazyReader(func() (io.Reader, error) {
		return strings.NewReader(in), nil
	})
	out, err := io.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestLazyReaderMulti(t *testing.T) {
	opened := 0
	r := io.MultiReader(
		NewLazyReader(func() (io.Reader, error) {
			opened += 1
			return strings.NewReader("hello world\n"), nil
		}),
		NewLazyReader(func() (io.Reader, error) {
			opened += 1
			return strings.NewReader("ciao planet"), nil
		}),
	)
	assert.Equal(t, 0, opened)
	out, err := io.ReadAll(r)
	assert.Equal(t, 2, opened)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\nciao planet", string(out))
}

func TestLazyReaderMultiGzip(t *testing.T) {
	opened := 0
	r, err := gzip.NewReader(io.MultiReader(
		NewLazyReader(func() (io.Reader, error) {
			opened += 1
			//
			buf := new(bytes.Buffer)
			gw := gzip.NewWriter(buf)
			if _, err := gw.Write([]byte("hello world\n")); err != nil {
				return nil, err
			}
			if err := gw.Flush(); err != nil {
				return nil, err
			}
			if err := gw.Close(); err != nil {
				return nil, err
			}
			data := buf.Bytes()
			//
			return bytes.NewReader(data), nil
		}),
		NewLazyReader(func() (io.Reader, error) {
			opened += 1
			//
			buf := new(bytes.Buffer)
			gw := gzip.NewWriter(buf)
			if _, err := gw.Write([]byte("ciao planet")); err != nil {
				return nil, err
			}
			if err := gw.Flush(); err != nil {
				return nil, err
			}
			if err := gw.Close(); err != nil {
				return nil, err
			}
			data := buf.Bytes()
			//
			return bytes.NewReader(data), nil
		}),
	))
	require.NoError(t, err)
	require.NotNil(t, r)
	// we expect the gzip reader to read the gzip header for the first file
	// when we call gzip.NewReader, so opened would be 1
	require.Equal(t, 1, opened)
	out, err := io.ReadAll(r)
	assert.Equal(t, 2, opened)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\nciao planet", string(out))
}

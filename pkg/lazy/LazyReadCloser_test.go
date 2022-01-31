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
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLazyReadCloser(t *testing.T) {
	in := "hello world"
	rc := NewLazyReadCloser(func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(in)), nil
	})
	assert.NoError(t, rc.Close()) // calling close before reading any data should return nil
	out, err := io.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
	assert.NoError(t, rc.Close())
}

func TestLazyReadCloserTwice(t *testing.T) {
	in := "hello world"
	rc := NewLazyReadCloser(func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(in)), nil
	})
	assert.NoError(t, rc.Close()) // calling close before reading any data should return nil
	out, err := io.ReadAll(rc)
	assert.NoError(t, err)
	assert.NoError(t, rc.Close())
	out, err = io.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
	assert.NoError(t, rc.Close())
}

func TestLazyReadCloserTwiceOnError(t *testing.T) {
	count := 0
	in := "hello world"
	rc := NewLazyReadCloser(func() (io.ReadCloser, error) {
		count++
		if count == 1 {
			return nil, fmt.Errorf("Fail on first count")
		}
		return io.NopCloser(strings.NewReader(in)), nil
	})
	assert.NoError(t, rc.Close()) // calling close before reading any data should return nil
	out, err := io.ReadAll(rc)
	assert.EqualError(t, err, "Fail on first count")
	assert.NoError(t, rc.Close())
	out, err = io.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
	assert.NoError(t, rc.Close())
}

func TestLazyReadCloserMulti(t *testing.T) {
	opened := 0
	r := io.MultiReader(
		NewLazyReadCloser(func() (io.ReadCloser, error) {
			opened += 1
			return io.NopCloser(strings.NewReader("hello world\n")), nil
		}),
		NewLazyReadCloser(func() (io.ReadCloser, error) {
			opened += 1
			return io.NopCloser(strings.NewReader("ciao planet")), nil
		}),
	)
	assert.Equal(t, 0, opened)
	out, err := io.ReadAll(r)
	assert.Equal(t, 2, opened)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\nciao planet", string(out))
}

func TestLazyReadCloserMultiGzip(t *testing.T) {
	opened := 0
	r, err := gzip.NewReader(io.MultiReader(
		NewLazyReadCloser(func() (io.ReadCloser, error) {
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
			return io.NopCloser(bytes.NewReader(data)), nil
		}),
		NewLazyReadCloser(func() (io.ReadCloser, error) {
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
			return io.NopCloser(bytes.NewReader(data)), nil
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

// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package lazy

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLazyFile(t *testing.T) {
	f := NewLazyFile("testdata/doc.txt", os.O_RDONLY, 0)
	out, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\n", string(out))
	fileInfo, err := f.Stat()
	assert.NoError(t, err)
	assert.Equal(t, "doc.txt", fileInfo.Name())
}

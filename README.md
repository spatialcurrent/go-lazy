[![CircleCI](https://circleci.com/gh/spatialcurrent/go-lazy/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-lazy/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-lazy)](https://goreportcard.com/report/spatialcurrent/go-lazy)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-lazy?status.svg)](https://godoc.org/github.com/spatialcurrent/go-lazy) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-lazy/blob/master/LICENSE)

# go-lazy

## Description

**go-lazy** is a library that includes a lazy reader and writer that allows delayed opening of a resource.  The `LazyReader` can be used to delay opening a file (and using a file descriptor) until the file is actually read.  The `LazyWriterAt` can be used to delay creating a file until the first byte is actually written.

# Usage

**Go**

You can import **go-lazy** as a library with:

```go
import (
  "github.com/spatialcurrent/go-lazy/pkg/lazy"
)
```

The easiest pattern is to simply use a closure function to open up a file or remote data source, so that you can delay opening the resource until it is actually needed.

```go
return lazy.NewReader(func() (io.Reader, error) {
  // open up a file or remote data source and return as a reader.
})
```

```go
writer := lazy.NewLazyWriterAt(func() (io.WriterAt, error) {
  // open up a file or byte buffer
})
```

A good use case is for reading multiple gzip files in sequential order.  The gzip reader can automatically read through multiple files.  This approach keeps the number of file descriptors open to 1 at a time while still using an io.MultiReader for convenience.

```
r, err := gzip.NewReader(io.MultiReader(
  NewLazyReader(func() (io.Reader, error) {
    return first file
  }),
  NewLazyReader(func() (io.Reader, error) {
    // return second file
  }),
))
```

Another good use case if for using the [s3manager](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager) Downloader to concurrently download multiple files.  To avoid many empty files for downloads that haven't started yet, use `LazyWriterAt`.

```
writer := lazy.NewLazyWriterAt(func() (io.WriterAt, error) {
  file, err := os.Create(input.Path)
  if err != nil {
    return nil, fmt.Errorf("error creating destination file %q: %w", input.Path, err)
  }
  return file, nil
})
_, err = downloader.DownloadWithContext(input.Context, writer, &s3.GetObjectInput{
  Bucket: aws.String(input.Bucket),
  Key:    aws.String(input.Key),
})
```

See [lazy](https://godoc.org/github.com/spatialcurrent/go-lazy/pkg/lazy) in GoDoc for API documentation.

# Testing

To run Go tests use `make test_go` (or `bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [ineffassign](https://github.com/gordonklaus/ineffassign), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-lazy/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.

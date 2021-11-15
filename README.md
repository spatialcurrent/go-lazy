[![CircleCI](https://circleci.com/gh/spatialcurrent/go-lazy/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-lazy/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-lazy)](https://goreportcard.com/report/spatialcurrent/go-lazy)  [![PkgGoDev](https://pkg.go.dev/badge/github.com/spatialcurrent/go-lazy)](https://pkg.go.dev/github.com/spatialcurrent/go-lazy) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-lazy/blob/master/LICENSE.md)

# go-lazy

## Description

**go-lazy** is a library that includes a lazy reader and writer that allows delayed opening of a resource.  The `LazyReader` and `LazyReadCloser` can be used to delay opening a resource until the resource is actually read.  The `LazyWriterAt` and `LazyWriteCloser` can be used to delay creating a resource until the first byte is actually written.  The `LazyFile` can be used to delay opening a file and allocating a file descriptor until it is required.

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
return lazy.NewReadCloser(func() (io.ReadCloser, error) {
  // open up a file or remote data source and return as a reader.
})
```

```go
writer := lazy.NewLazyWriteCloser(func() (io.WriterAt, error) {
  // open up a pipe
})
```

```go
writer := lazy.NewLazyWriterAt(func() (io.WriterAt, error) {
  // open up a file or byte buffer
})
```

A good use case is for reading multiple gzip files in sequential order.  The gzip reader can automatically read through multiple files.  This approach keeps the number of file descriptors open to 1 at a time while still using an io.MultiReader for convenience.

```go
r, err := gzip.NewReader(io.MultiReader(
  NewLazyReader(func() (io.Reader, error) {
    return first file
  }),
  NewLazyReader(func() (io.Reader, error) {
    // return second file
  }),
))
```

The LazyFile can also be used to delay opening or stating a file if it is passed to a callback.  LazyFile matches the `os.File` API when possible, but is not a drop in replacement.  The API is identical for `Read`, `Stat`, and `Close` methods, but the `Fd` method now returns the file descriptor with an error if any.

```go
f := lazy.NewLazyFile(name, os.O_RDONLY, 0)
err := callback(f)
if err != nil {
  return err
}
// do something with the file
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

See [go.dev](https://pkg.go.dev/github.com/spatialcurrent/go-lazy/) for information on how to use Go API.

# Testing

To run Go tests use `make test_go` (or `bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [ineffassign](https://github.com/gordonklaus/ineffassign), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-lazy/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.

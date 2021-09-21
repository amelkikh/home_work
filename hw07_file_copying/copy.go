package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeValue         = errors.New("negative value is unacceptable")
)

type ProgressReader struct {
	io.Reader

	total int64
}

func (pt *ProgressReader) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if err != nil {
		return 0, err
	}
	pt.total += int64(n)

	return n, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	info, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if limit < 0 || offset < 0 {
		return ErrNegativeValue
	}
	if info.Size() == 0 {
		return ErrUnsupportedFile
	}
	if offset > info.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = info.Size() - offset
	}

	f, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fTo.Close()

	_, err = f.Seek(offset, 0)
	if err != nil {
		return err
	}
	bar := pb.New64(limit)
	bar.Start()
	defer bar.Finish()

	r := bar.NewProxyReader(f)
	w := bufio.NewWriter(fTo)

	_, err = io.CopyN(w, r, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

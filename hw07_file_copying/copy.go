package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeValue         = errors.New("negative value is unacceptable")
)

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

	r := bufio.NewReader(f)
	if offset > 0 {
		_, err = r.Discard(int(offset))
		if err != nil {
			return err
		}
	}
	w := bufio.NewWriter(fTo)
	buf := make([]byte, 1024)

	bar := pb.StartNew(int(limit))
	bar.Set(pb.Bytes, true)
	defer bar.Finish()

	var n int64 = 0
	for n <= limit {
		nn, err := r.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		if nn == 0 {
			break
		}

		if n+int64(nn) > limit {
			nn = int(limit - n)
			n = limit
		} else {
			n += int64(nn)
		}

		if _, err := w.Write(buf[:nn]); err != nil {
			return err
		}

		bar.Add(nn)
		// Just for testing purposes to visualize progress
		time.Sleep(time.Millisecond * 300)
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

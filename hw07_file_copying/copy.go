package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWritingToFile         = errors.New("error while writing to a file")
	ErrUnexpectedEOF         = errors.New("unexpected EOF")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		offset = 0
	}

	if limit < 0 {
		limit = 0
	}

	state, err := os.Stat(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	if state.IsDir() {
		return ErrUnsupportedFile
	}

	size := state.Size()
	if offset >= size {
		return ErrOffsetExceedsFileSize
	}

	amount := size

	if limit != 0 && limit < size {
		amount = limit
	}

	if offset != 0 {
		if limit == 0 {
			amount = size - offset
		} else if offset+limit >= size {
			amount = size - offset
		}
	}

	fileFrom, errF := os.Open(fromPath)

	if errF != nil {
		return errF
	}

	defer func() {
		_ = fileFrom.Close()
	}()

	if offset > 0 {
		_, err := fileFrom.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	fileTo, errT := os.Create(toPath)

	if errT != nil {
		return errT
	}

	defer func() {
		_ = fileTo.Close()
	}()

	progressBar := pb.Start64(amount)
	barReader := progressBar.NewProxyReader(fileFrom)
	progressBar.Start()

	number, errN := io.CopyN(fileTo, barReader, amount)

	if errN != nil {
		if errors.Is(errN, io.EOF) {
			return ErrUnexpectedEOF
		}

		return errN
	}

	progressBar.Finish()

	if number != amount {
		return ErrWritingToFile
	}

	return nil
}

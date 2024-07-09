package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile         = errors.New("unsupported file")
	ErrOffsetExceedsFileSize   = errors.New("offset exceeds file size")
	ErrSeekerValueMoreThanZero = errors.New("seeker values must be positive")
)

func checkFromFile(offset, limit, fileSize int64) error {
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	if limit < 0 || offset < 0 {
		return ErrSeekerValueMoreThanZero
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// получаем файл-источник
	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0o666)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		panic(err)
	}
	defer file.Close()
	// получаем информацию о файле-источнике
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	err = checkFromFile(offset, limit, fileInfo.Size())
	if err != nil {
		return err
	}

	// прочитать файл-источник подготовить reader

	// готовим writer
	dst, err := os.CreateTemp("", "tmp")
	if err != nil {
		return err
	}

	_, err = io.CopyN(dst, file, limit)
	if err != nil {
		return err
	}

	return nil
}

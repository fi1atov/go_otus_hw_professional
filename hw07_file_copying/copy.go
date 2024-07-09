package main

import (
	"errors"
	"fmt"
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
	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
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

	// готовим writer
	dst, err := os.CreateTemp(toPath, "tmp")
	if err != nil {
		return err
	}
	defer dst.Close()

	fmt.Println("file: ", file)
	fmt.Println("dst: ", dst)

	if limit == 0 {
		limit = fileInfo.Size()
	}

	value, err := io.CopyN(dst, file, limit)
	if err != nil {
		return err
	}
	fmt.Println(value)

	return nil
}

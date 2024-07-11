package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
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

func Copy(fromPath, toPath string, offset, limit int64) (string, error) {
	// получаем файл-источник
	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return "", err
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
		return "", err
	}

	// двигаем указатель
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return "", err
	}

	// готовим writer
	// получаем файл-назначение на чтение/запись + нужна очистка файла перед его заполнением
	dst, err := os.OpenFile(toPath, os.O_RDWR|os.O_TRUNC, 0o666)
	if err != nil {
		if os.IsNotExist(err) {
			dst, err = os.CreateTemp("mydir", toPath)
			if err != nil {
				return "", err
			}
		}
	}
	defer dst.Close()

	// для того чтобы при нуле читали все и чтобы избежать ошибку EOF
	// при слишком большом лимите
	if limit == 0 || limit > fileInfo.Size() {
		limit = fileInfo.Size()
	}

	// для того чтобы избежать ошибку EOF
	if offset != 0 && offset+limit > fileInfo.Size() {
		limit = fileInfo.Size() - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(file)

	_, err = io.CopyN(dst, barReader, limit)

	bar.Finish()

	if err != nil {
		return "", err
	}

	return dst.Name(), nil
}

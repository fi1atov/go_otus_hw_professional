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

func Copy(fromPath, toPath string, offset, limit int64) error {
	var dst *os.File
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

	// двигаем указатель
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// для того чтобы при нуле читали все и чтобы избежать ошибку EOF
	// при слишком большом лимите
	if limit == 0 || limit > fileInfo.Size() {
		limit = fileInfo.Size()
	}

	// для того чтобы избежать ошибку EOF
	if offset != 0 && offset+limit > fileInfo.Size() {
		limit = fileInfo.Size() - offset
	}

	// источник и назначение совпадают?
	if fromPath == toPath {
		// совпадают - надо создавать временный файл и копировать него
		// вторым шагом - удаляем src-файл и создаем файл заново и наполняем его
		// получим путь до текущей директории
		path, err := os.Getwd()
		if err != nil {
			return err
		}
		// готовим writer
		// создаем временный файл
		dst, err = os.CreateTemp(path, "*"+toPath)
		if err != nil {
			return err
		}
		defer dst.Close()
	} else {
		// не совпадает - временный файл не нужен - создаем сразу целевой файл
		// готовим writer
		// получаем файл-назначение на чтение/запись + нужна очистка файла перед его заполнением
		dst, err = os.Create(toPath)
		if err != nil {
			return err
		}
		defer dst.Close()
	}
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(file)

	_, err = io.CopyN(dst, barReader, limit)
	if err != nil {
		return err
	}

	bar.Finish()

	if fromPath == toPath {
		err := os.Remove(fromPath)
		if err != nil {
			return err
		}
		err = os.Rename(dst.Name(), toPath)
		if err != nil {
			return err
		}
	}

	return nil
}

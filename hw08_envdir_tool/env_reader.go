package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadEnvDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)
	for _, file := range files {
		// получаем информацию о файле-источнике
		fileInfo, err := file.Info()
		if err != nil {
			panic(err)
		}
		// имя файла не должно содержать '=' - игнорируем такие файлы
		if strings.Contains(file.Name(), "=") {
			continue
		}
		// если файл полностью пустой - удаляем переменную окружения
		if fileInfo.Size() == 0 {
			result[file.Name()] = EnvValue{NeedRemove: true}
			continue
		}
		// проверки пройдены - открываем, читаем файл, получаем значение переменной окружения
		value, err := readValue(dir, file.Name())
		if err != nil {
			return nil, err
		}
		result[file.Name()] = EnvValue{Value: value}
	}

	return result, nil
}

func readValue(dir, fileName string) (string, error) {
	f, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return "", err
	}
	// закрыть файл и обработать возможную ошибку при закрытии
	defer func() {
		mustNil(f.Close())
	}()

	s := bufio.NewScanner(f)
	if !s.Scan() {
		return "", nil
	}
	line := s.Text()
	// терминальные нули (0x00) заменяются на перевод строки
	line = strings.ReplaceAll(line, "\x00", "\n")
	// пробелы и табуляция в конце удаляются
	line = strings.TrimRightFunc(line, unicode.IsSpace)
	return line, nil
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

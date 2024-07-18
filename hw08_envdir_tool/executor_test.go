package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	capturer "github.com/zenizh/go-capturer"
)

func TestRunCmd(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// подготовка тестовых данных - создаем новую временную папку
		dir, err := os.MkdirTemp("", "test")
		require.NoError(t, err) // при создании папки ошибки быть не должно
		defer os.RemoveAll(dir) // в конце теста удалим созданную на время теста папку

		// папка с переменными окружения
		err = os.Mkdir(filepath.Join(dir, "vars"), 0o777) // создадим папку vars внутри временной папки для тестов
		require.NoError(t, err)                           // при создании папки ошибки быть не должно
		// создадим файл "BAR" и запишем в него "bar"
		err = os.WriteFile(filepath.Join(dir, "vars", "BAR"), []byte("bar"), 0o666)
		require.NoError(t, err)
		// создади файл - bash-скрипт, внутрь него поместим bash-команду (упрощенную версию для юнит-теста)
		// которая распечатывает свой первый аргумент который мы можем передать и переменную окружения BAR
		err = os.WriteFile(filepath.Join(dir, "t.sh"), []byte("#!/usr/bin/env bash\necho $1\necho $BAR\n"), 0o666)
		require.NoError(t, err)
		err = os.Chmod(filepath.Join(dir, "t.sh"), 0o777) // настроили execute/read/write доступ для созданного скрипта
		require.NoError(t, err)
		// конец подготовки тестовых данных

		env, err := ReadEnvDir(filepath.Join(dir, "vars")) // читаем созданную в тесте переменную из созданной в тесте папки
		require.NoError(t, err)

		var returnCode int
		// приходится получать результат программы именно из stout, поэтому использую go-capturer
		result := capturer.CaptureStdout(func() {
			// запускаем подготовленный тестовый скрипт - передаем аргумент + переменную окружения
			returnCode = RunCmd([]string{filepath.Join(dir, "t.sh"), "something"}, env)
		})
		require.Equal(t, 0, returnCode)
		// вывод должен соответствовать ожиданию (выводится переданный аргумент + переменная окружения)
		require.Equal(t, "something\nbar\n", result)
	})

	t.Run("Rewrite FOO", func(t *testing.T) {
		// подготовка тестовых данных
		dir, err := os.MkdirTemp("", "test")
		require.NoError(t, err) // при создании папки ошибки быть не должно
		defer os.RemoveAll(dir) // в конце теста удалим созданную на время теста папку

		// папка с переменными окружения
		err = os.Mkdir(filepath.Join(dir, "vars"), 0o777) // создадим папку vars внутри временной папки для тестов
		require.NoError(t, err)                           // при создании папки ошибки быть не должно
		// создадим файл "FOO" и запишем в него "42"
		err = os.WriteFile(filepath.Join(dir, "vars", "FOO"), []byte("42"), 0o666)
		require.NoError(t, err)
		// создади файл - bash-скрипт, внутрь него поместим bash-команду (упрощенную версию для юнит-теста)
		// которая распечатывает переменную окружения FOO
		err = os.WriteFile(filepath.Join(dir, "t.sh"), []byte("#!/usr/bin/env bash\necho $FOO\n"), 0o666)
		require.NoError(t, err)
		err = os.Chmod(filepath.Join(dir, "t.sh"), 0o777) // настроили execute/read/write доступ для созданного скрипта
		require.NoError(t, err)
		// конец подготовки тестовых данных

		env, err := ReadEnvDir(filepath.Join(dir, "vars")) // читаем созданную в тесте переменную из созданной в тесте папки
		require.NoError(t, err)

		var returnCode int
		result := capturer.CaptureStdout(func() {
			returnCode = RunCmd([]string{filepath.Join(dir, "t.sh")}, env)
		})
		require.Equal(t, 0, returnCode)
		// вывод должен соответствовать ожиданию (значение переменной окружения равно значению из файла)
		require.Equal(t, "42\n", result)
	})
}

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Success with testdata", func(t *testing.T) {
		expectEnv := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{Value: ""},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: `"hello"`},
			"UNSET": EnvValue{NeedRemove: true},
		}
		env, err := ReadEnvDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, env, expectEnv)
	})

	t.Run("Success", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		// файл, который должен быть проигнорирован
		err = os.WriteFile(filepath.Join(dir, "t=t"), []byte("bar"), 0o666)
		require.NoError(t, err)
		// имя файла маленькими буквами
		err = os.WriteFile(filepath.Join(dir, "test"), []byte("test"), 0o666)
		require.NoError(t, err)
		// файл с пустой первой строкой
		err = os.WriteFile(filepath.Join(dir, "EMPTY"), []byte("\n"), 0o666)
		require.NoError(t, err)

		expectEnv := Environment{
			"test":  EnvValue{Value: "test", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
		}
		env, err := ReadEnvDir(dir)
		require.NoError(t, err)
		require.Equal(t, env, expectEnv)
	})

	t.Run("Success with empty dir", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		env, err := ReadEnvDir(dir)
		require.NoError(t, err)
		require.Len(t, env, 0)
	})

	t.Run("Fail with dir not exists", func(t *testing.T) {
		_, err := ReadEnvDir("some name")
		require.Error(t, err)
	})
}

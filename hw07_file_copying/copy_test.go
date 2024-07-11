package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		desc             string
		fromPath, toPath string
		offset, limit    int64
		expectedFile     string
	}{
		{
			desc:         "offset0_limit0",
			fromPath:     "testdata/input.txt",
			toPath:       "out",
			offset:       0,
			limit:        0,
			expectedFile: "testdata/out_offset0_limit0.txt",
		},
		{
			desc:         "offset0_limit10",
			fromPath:     "testdata/input.txt",
			toPath:       "out",
			offset:       0,
			limit:        10,
			expectedFile: "testdata/out_offset0_limit10.txt",
		},
		{
			desc:         "offset0_limit1000",
			fromPath:     "testdata/input.txt",
			toPath:       "out",
			offset:       0,
			limit:        1000,
			expectedFile: "testdata/out_offset0_limit1000.txt",
		},
		{
			desc:         "offset0_limit10000",
			fromPath:     "testdata/input.txt",
			toPath:       "out",
			offset:       0,
			limit:        10000,
			expectedFile: "testdata/out_offset0_limit10000.txt",
		},
		{
			desc:         "offset100_limit1000",
			fromPath:     "testdata/input.txt",
			toPath:       "out",
			offset:       100,
			limit:        1000,
			expectedFile: "testdata/out_offset100_limit1000.txt",
		},
		{
			desc:         "offset6000_limit1000",
			fromPath:     "testdata/input.txt",
			toPath:       "out",
			offset:       6000,
			limit:        1000,
			expectedFile: "testdata/out_offset6000_limit1000.txt",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// выполняем копирование
			dstFile, err := Copy(tc.fromPath, tc.toPath, tc.offset, tc.limit)
			// ошибок быть не должно
			require.NoError(t, err)
			// смотрим что там скопировалось и сверим с нашим ожиданием
			// откроем первый файл
			fileFirst, err := os.OpenFile(dstFile, os.O_RDONLY, 0)
			if err != nil {
				panic(err)
			}
			defer fileFirst.Close()
			// получаем информацию о файле который был создан во время тестирования
			fileFirstInfo, err := fileFirst.Stat()
			if err != nil {
				panic(err)
			}

			// откроем второй файл
			fileSecond, err := os.OpenFile(tc.expectedFile, os.O_RDONLY, 0)
			if err != nil {
				panic(err)
			}
			defer fileSecond.Close()
			// получаем информацию о файле который мы ожидаем
			fileSecondInfo, err := fileSecond.Stat()
			if err != nil {
				panic(err)
			}

			// Сверка размеров созданного  во время теста файла и ожидаемого файла
			require.Equal(t, fileSecondInfo.Size(), fileFirstInfo.Size())

			// Удалим созданный во время теста файл
			// (при необходимости просматривать файлы - закомментировать эту строку удаления)
			os.Remove(dstFile)
		})
	}
}

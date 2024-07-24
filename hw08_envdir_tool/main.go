package main

import (
	"fmt"
	"os"
)

// envdir возвращает 111 в случае если возникает ошибка чтения файлов из директории.
const selfReturnCode = 111

func main() {
	// прочитать аргументы начиная со второго (чтобы не хватать путь к директории)
	args := os.Args[1:]
	// открываем директорию в которой находятся файлы-переменные окружения
	env, err := ReadEnvDir(args[0])
	if err != nil {
		printFatal(fmt.Errorf("go-envdir: ошибка: %w", err))
	}
	code := RunCmd(args[1:], env)
	if code == selfReturnCode {
		printFatal(fmt.Errorf("go-envdir: программа не запускается %s", args[1]))
	}
	os.Exit(code)
}

func printFatal(a interface{}) {
	// log.Fatal не подходит, потому что добавляет перед сообщением дополнительную информацию,
	// а для утилиты это неправильное поведение
	_, _ = fmt.Fprintln(os.Stderr, a)
	os.Exit(selfReturnCode)
}

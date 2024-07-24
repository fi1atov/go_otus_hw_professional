package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)
	// передать в программу слайс переменных окружения
	command.Env = changeEnv(env)
	// стандартные потоки ввода/вывода/ошибок нужно пробросить в вызываемую программу
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return selfReturnCode
	}
	return 0
}

func changeEnv(env Environment) []string {
	result := os.Environ()
	// смотрим файлы-переменные окружения которые получили из директории
	for name, v := range env {
		result = filterEnv(result, name)
		// если ее не нужно исключить - добавляем в слайс переменных
		if !v.NeedRemove {
			result = append(result, name+"="+v.Value)
		}
	}
	return result
}

func filterEnv(envVars []string, delName string) []string {
	n := 0
	for _, v := range envVars {
		// перебираем все имеющиеся в системе переменные окружения - извлекаю имена
		name := strings.SplitN(v, "=", 2)[0]
		// добавляем недостающую переменную окружения из файла в общий список переменных
		if name != delName {
			envVars[n] = v
			n++
		}
	}
	return envVars[:n]
}

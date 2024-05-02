package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func Unpack(s string) (r string, err error) {
	// если строка по каким-то причинам не превращается в числа
	if _, err := strconv.Atoi(s); err == nil {
		return r, ErrInvalidString
	}
	// если пришла пустая строка - возвращаем ее
	if s == "" {
		return "", err
	}
	// если первый символ строки - это число
	if char := rune(s[0]); unicode.IsDigit(char) {
		return r, ErrInvalidString
	}

	var prev rune
	var prevPrevChar rune
	var escaped bool
	var ecran bool
	var b strings.Builder
	for _, char := range s {
		// если предыдущий и текущий символы являются цифрами и нет экранирования - ошибка
		if string(prevPrevChar) != "\\" && unicode.IsDigit(prev) && unicode.IsDigit(char) {
			return r, ErrInvalidString
		}
		// если текущий символ не цифра и не слеш, и при этом есть экранирование - ошибка
		if !unicode.IsDigit(char) && string(char) != "\\" && ecran {
			return r, ErrInvalidString
		}

		switch unicode.IsDigit(char) && !escaped {
		case true:
			m := int(char - '0') // удобный способ превращения руны в число
			if m == 0 {
				// если пришел ноль - значит букву перед этим нулем нужно убрать
				res := b.String()                // получаем строку из формируемого буфера
				b.Reset()                        // очистить буфер
				b.WriteString(trimLastChar(res)) // убираем последний символ и перезаписываем буфер
			} else {
				r := strings.Repeat(string(prev), m-1) // повторить символ на пришедшее число - 1, т.к. 1 символ уже вписан
				b.WriteString(r)
				if ecran {
					ecran = false // экранированный символ записан - снять экранирование
				}
			}
		case false:
			escaped = string(char) == "\\" && !ecran
			if escaped {
				ecran = true // пришел символ экранирования
			} else {
				b.WriteRune(char)
				ecran = false // экранирование тут 100% надо снимать
			}
		}
		// запоминать один предыдущий символ недостаточно - нужно запоминать два предущих символа
		prevPrevChar = prev // предпредыдущий символ
		prev = char         // предыдущий литерал
	}

	// если последний символ строки - это символ экранирования
	if ecran {
		return r, ErrInvalidString
	}

	return b.String(), err
}

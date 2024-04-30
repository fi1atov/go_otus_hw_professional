package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

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
	var prevChar rune
	var escaped bool
	var ecran bool
	var b strings.Builder
	for _, char := range s {
		fmt.Println(char, string(char))
		// fmt.Println(escaped)
		// предыдущий и текущий символы строки могут быть числами только если было экранирование
		if unicode.IsDigit(prevChar) && unicode.IsDigit(char) && !ecran {
			return r, ErrInvalidString
		}
		// если текущий символ не цифра и не слеш, и при этом есть экранирование - ошибка
		if !unicode.IsDigit(char) && string(char) != "\\" && ecran {
			return r, ErrInvalidString
		}

		m := int(char - '0') // удобный способ превращения руны в число
		switch unicode.IsDigit(char) && !escaped {
		case true:
			if m == 0 {
				// если пришел ноль - значит букву перед этим нулем нужно убрать
				res := b.String()               // получаем строку из формируемого буфера
				b.Reset()                       // очистить буфер
				b.WriteString(res[:len(res)-1]) // убираем последний символ и перезаписываем буфер
			} else {
				r := strings.Repeat(string(prev), m-1) // повторить символ на пришедшее число - 1, т.к. 1 символ уже вписан
				b.WriteString(r)
				if ecran {
					ecran = false // экранированный символ записан - снять экранирование
				}
			}
		case false:
			escaped = string(char) == "\\" && string(prev) != "\\"
			if !escaped {
				b.WriteRune(char)
				// if ecran {
				// 	ecran = false // экранированный символ записан - снять экранирование
				// }
			} else {
				ecran = true // пришел символ экранирования
			}
			prev = char
		}
		prevChar = char
	}

	return b.String(), err
}

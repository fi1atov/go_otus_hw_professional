package main

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
	var b strings.Builder
	for _, char := range s {
		// если предыдущий и текущий символы строки - это числа
		if unicode.IsDigit(prevChar) && unicode.IsDigit(char) {
			return r, ErrInvalidString
		}
		if unicode.IsDigit(char) && !escaped {
			m := int(char - '0') // превратить в число
			// fmt.Println("m", m)
			if m == 0 {
				res := b.String()
				b.Reset()
				b.WriteString(res[:len(res)-1])
				continue
			}
			r := strings.Repeat(string(prev), m-1)
			// fmt.Println("r", r)
			b.WriteString(r)
		} else {
			escaped = string(char) == "\\" && string(prev) != "\\"
			if !escaped {
				// fmt.Println(char)
				b.WriteRune(char)
			}
			prev = char
		}
		prevChar = char
	}

	return b.String(), err
}

func main() {
	str := "d\n5abc"
	fmt.Println(Unpack(str))
}

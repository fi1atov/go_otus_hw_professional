package hw10programoptimization

// отказываемся от regexp потому что он медленный и используем быстрый jsoniter
// используем bufio для получения пользователей.
import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u := getUsers(r)
	return countDomains(u, domain)
}

type NextUser = func() (*User, bool, error)

func getUsers(r io.Reader) NextUser {
	// открыли reader и отдаем функцию которая может из него читать пока не появится EOF
	var auser User
	s := bufio.NewScanner(r)

	nextUser := func() (user *User, ok bool, err error) {
		// чтение с помощью Scan() и Bytes()
		ok = s.Scan()

		if !ok {
			err = s.Err()
			if err != nil {
				err = fmt.Errorf("ошибка сканирования: %w", err)
			}
			return
		}

		if err = jsoniter.Unmarshal(s.Bytes(), &auser); err != nil {
			err = fmt.Errorf("ошибка чтения пользователя: %w", err)
			return
		}

		user = &auser
		// fmt.Println(user)
		return
	}

	return nextUser
}

func countDomains(nextUser NextUser, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainMask := "." + domain

	for {
		// читаем строку, курсор будет остановлен в конце каждого прочитанного объекта
		user, ok, err := nextUser()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		email := strings.ToLower(user.Email)
		// email оканчивается на наш домен
		matched := strings.HasSuffix(email, domainMask)

		if matched {
			n := strings.LastIndex(email, "@") // возвращает -1 if substr is not present in s
			if n == -1 {
				return nil, fmt.Errorf("email не совпал: %s", email)
			}
			result[email[n+1:]]++ // считаем совпадения
		}
	}

	return result, nil
}

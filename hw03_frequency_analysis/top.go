package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Word struct {
	Word  string
	Count int
}

var sortedStruct []Word

func separateRow(input string) []string {
	return regexp.MustCompile(`(?:\s-\s)|[,.\\s]+`).Split(input, -1)
}

func Top10(input string) []string {
	// 1. Отделить слова и превратить текст в слайс
	dict := make(map[string]int)
	words := separateRow(input)
	wordsStr := strings.ToLower(strings.Join(words, " "))
	// Избавляет от отступов, переноса строки, пробелов
	words = strings.Fields(wordsStr)
	// Благодаря свойству мап(неповторимость ключа)
	// собираеми и считаем слова
	for _, word := range words {
		dict[word]++
	}
	// мапа не имеет функции сортировки - нужно преобразовать ее в слайс структур
	for key, value := range dict {
		sortedStruct = append(sortedStruct, Word{key, value})
	}
	// слайс структур теперь можно отсортировать по значению и лексикографически
	sort.Slice(sortedStruct, func(i, j int) bool {
		if sortedStruct[i].Count != sortedStruct[j].Count {
			return sortedStruct[i].Count > sortedStruct[j].Count // по значению
		}
		return sortedStruct[i].Word < sortedStruct[j].Word // лексикографически
	})
	// слайс структур должен стать слайсом текстов - интересуют первые 10 слов
	var result []string
	for i, p := range sortedStruct {
		if i < 10 {
			result = append(result, p.Word)
		} else {
			break
		}
	}
	return result
}

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

// Скомпилированное регулярное выражение на уровне пакета.
var wordSplitRegex = regexp.MustCompile(`(?:\s-\s)|[,.\\s]+`)

func separateRow(input string) []string {
	return wordSplitRegex.Split(input, -1)
}

func CountWords(input string) map[string]int {
	dict := make(map[string]int)
	// Отделить слова и превратить текст в слайс
	words := separateRow(input)
	wordsStr := strings.ToLower(strings.Join(words, " "))
	// Избавляет от отступов, переноса строки, пробелов
	words = strings.Fields(wordsStr)
	// Благодаря свойству мап(неповторимость ключа)
	// собираеми и считаем слова
	for _, word := range words {
		dict[word]++
	}
	return dict
}

func Top10(input string) []string {
	// посчитать количество повторений слова в тексте
	dict := CountWords(input)
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

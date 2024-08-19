package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false
var re = regexp.MustCompile(`^[^\w\sа-яА-ЯёЁ]+|[^\w\sа-яА-ЯёЁ]+$`)

type wordInfo struct {
	word  string
	count int
}

func Top10(text string) []string {
	words := strings.Fields(text)
	if len(words) < 2 {
		return words
	}

	wordFrequencies := make(map[string]int)

	for _, word := range words {
		if taskWithAsteriskIsCompleted {
			word = re.ReplaceAllString(strings.ToLower(word), "")
			if len(word) < 1 {
				continue
			}

		}

		if count, found := wordFrequencies[word]; found {
			wordFrequencies[word] = count + 1
			continue
		}

		wordFrequencies[word] = 1
	}

	wordInfos := make([]wordInfo, 0, len(wordFrequencies))
	for word, count := range wordFrequencies {
		wordInfos = append(wordInfos, wordInfo{word: word, count: count})
	}

	sort.Slice(wordInfos, func(i, j int) bool {
		if wordInfos[i].count == wordInfos[j].count {
			return wordInfos[i].word < wordInfos[j].word
		}
		return wordInfos[i].count > wordInfos[j].count
	})

	result := make([]string, 0, 10)
	for _, wordInfo := range wordInfos {
		result = append(result, wordInfo.word)
	}

	if len(result) < 11 {
		return result
	}

	return result[:10]
}

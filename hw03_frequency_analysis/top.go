package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	re       = regexp.MustCompile(`(?mi)([\p{L}-]+)`)
	excludes = map[string]struct{}{
		"-": {},
	}
	wordsLimit = 10
)

func Top10(str string) []string {
	str = strings.ToLower(str)
	m := make(map[string]int)
	for _, match := range re.FindAllString(str, -1) {
		if _, ok := excludes[match]; ok {
			continue
		}
		m[match]++
	}

	if len(m) == 0 {
		return []string{}
	}

	freqList := make(WordFrequencies, 0, len(m))
	for word, count := range m {
		freqList = append(freqList, Pair{
			Word:  word,
			Count: count,
		})
	}
	sort.Slice(freqList, func(i, j int) bool {
		if a, b := freqList[i].Count, freqList[j].Count; a != b {
			return a > b
		}
		return freqList[i].Word < freqList[j].Word
	})

	if wordsLimit > len(freqList) {
		wordsLimit = len(freqList)
	}

	result := make([]string, 0, wordsLimit)
	for i := 0; i < wordsLimit; i++ {
		result = append(result, freqList[i].Word)
	}

	return result
}

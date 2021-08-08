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
		if _, ok := m[match]; ok {
			m[match]++
			continue
		} else {
			m[match] = 1
		}
	}

	if len(m) == 0 {
		return []string{}
	}

	freqList := make(WordFrequencies, len(m))
	i := 0
	for word, count := range m {
		freqList[i] = Pair{
			Word:  word,
			Count: count,
		}
		i++
	}
	sort.Sort(sort.Reverse(freqList))

	// Split and sort alphabetically each equal chunk of counters
	currentCounts := 0
	result := make([]string, 0)
	subSlice := make([]string, 0)
	for _, pair := range freqList {
		if pair.Count != currentCounts {
			currentCounts = pair.Count
			if len(subSlice) > 0 {
				sort.Strings(subSlice)
				result = append(result, subSlice...)
				subSlice = make([]string, 0)
			}
		}
		subSlice = append(subSlice, pair.Word)
	}
	sort.Strings(subSlice)
	result = append(result, subSlice...)

	if wordsLimit > len(result) {
		wordsLimit = len(result)
	}
	return result[:wordsLimit]
}

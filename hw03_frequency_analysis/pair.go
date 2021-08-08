package hw03frequencyanalysis

type Pair struct {
	Word  string
	Count int
}

type WordFrequencies []Pair

func (p WordFrequencies) Len() int           { return len(p) }
func (p WordFrequencies) Less(i, j int) bool { return p[i].Count < p[j].Count }
func (p WordFrequencies) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type AlphabeticalAscending []Pair

func (p AlphabeticalAscending) Len() int           { return len(p) }
func (p AlphabeticalAscending) Less(i, j int) bool { return p[i].Word < p[j].Word }
func (p AlphabeticalAscending) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

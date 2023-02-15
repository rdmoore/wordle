package words

const (
	wordLength = 5
)

type RuneCounts map[rune][]int

func RuneOccurrences(words []*Word) RuneCounts {
	runeCounts := make(RuneCounts)
	for _, w := range words {
		for i, ch := range w.Text {
			if runeCounts[ch] == nil {
				runeCounts[ch] = make([]int, wordLength)
			}
			runeCounts[ch][i] += 1
		}
	}
	return runeCounts
	// tuples := make([]letterTuple, 0, len(runeCounts))
	// for l, c := range runeCounts {
	// 	tuples = append(tuples, letterTuple{
	// 		letter: l,
	// 		count:  c,
	// 	})
	// }
	// sort.Slice(tuples, func(i, j int) bool {
	// 	if tuples[i].count == tuples[j].count {
	// 		return tuples[i].letter > tuples[j].letter
	// 	}
	// 	return tuples[i].count > tuples[j].count
	// })
	// return tuples
}

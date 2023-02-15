package command

import (
	"fmt"
	"strings"

	"wordle/words"

	"github.com/sirupsen/logrus"
)

// const (
// 	vowels = "aeiou"
// )

type rankCommand struct {
	dictionary    string
	excluded      string
	includes      []string
	exactPosition string
	targets       []string
	// vowels        bool
}

func NewRank(app Commander) (string, RunFunc) {
	var args rankCommand
	cmd := app.Command("rank", "rank words based on strategy")
	cmd.Flag("dictionary", "file from which to read dictionary words").Default("five-letter-words.txt").StringVar(&args.dictionary)
	cmd.Flag("include", "leters that must be present, but not in position specified").Short('i').StringsVar(&args.includes)
	cmd.Flag("exclude", "leters that are not part of the target word").Short('e').StringVar(&args.excluded)
	cmd.Flag("exact", "word containing letters in exact positions using space for unknown. Example ' a  e'").Short('x').StringVar(&args.exactPosition)
	cmd.Flag("target", "words for which verbose logging should be enabled").Short('t').StringsVar(&args.targets)
	// cmd.Flag("vowels", "rank vowels more important").BoolVar(&args.vowels)
	return cmd.FullCommand(), args.Run
}

func (cmd *rankCommand) Run() error {
	// TODO: need a better approach to setting this - a way to select the logger based on either the word
	// or a generate verbose flag.
	logrus.SetLevel(logrus.DebugLevel)

	if err := cmd.validate(); err != nil {
		return err
	}
	matched, err := words.LoadWords(cmd.dictionary, cmd.filterWord)
	if err != nil {
		return err
	}
	if len(matched) == 0 {
		return fmt.Errorf("no words selected from input")
	}
	weights := words.RuneOccurrences(matched)
	scores := words.Scores(matched, weights)

	if len(scores) == 0 {
		fmt.Printf("no words available after scoring\n")
		return nil
	}
	best := scores[0].Total
	for i, tuple := range scores {
		if scores[i].Total < best && i > 30 {
			break
		}
		fmt.Printf("%4d - %s\n", i, tuple.String())
	}
	index := len(scores) - 1
	if scores[index].Total != best && index > 30 {
		fmt.Printf("lowest ranked word: %4d - %s\n", index, scores[index].String())
	}
	return nil
}

func (cmd *rankCommand) validate() error {
	var targets []string
	for _, w := range cmd.targets {
		targets = append(targets, strings.Split(w, ",")...)
	}
	cmd.targets = targets
	for _, word := range cmd.targets {
		if !words.ValidWord.MatchString(word) {
			return fmt.Errorf("invalid target word [%s]", word)
		}
	}
	if cmd.excluded != "" && !words.ValidExclude.MatchString(cmd.excluded) {
		return fmt.Errorf("exclude string [%s] contains invalid character(s)", cmd.excluded)
	}
	if cmd.exactPosition != "" && !words.ValidInclude.MatchString(cmd.exactPosition) {
		return fmt.Errorf("exact [%s] contains invalid character(s)", cmd.exactPosition)
	}
	if strings.ContainsAny(cmd.exactPosition, cmd.excluded) {
		return fmt.Errorf("exact [%s] contains excluded character(s)", cmd.exactPosition)
	}
	runes := []rune(cmd.exactPosition)
	for _, word := range cmd.includes {
		if !words.ValidInclude.MatchString(word) {
			return fmt.Errorf("include [%s] contains invalid character(s)", word)
		}
		if len(runes) > 0 {
			for i, r := range word {
				switch r {
				case '.', ' ': // allowable "nil" characters
					continue
				case runes[i]:
					return fmt.Errorf("exact word [%s] overlaps included [%s] at character [%c]", cmd.exactPosition, word, r)
				default:
				}
			}
		}
	}
	return nil
}

func (cmd *rankCommand) verbose(word *words.Word) bool {
	for _, w := range cmd.targets {
		if word.Text == w {
			return true
		}
	}
	return false
}

func (cmd *rankCommand) filterWord(word *words.Word) bool {
	verbose := cmd.verbose(word)
	if words.Excluded(word, cmd.excluded, verbose) {
		if verbose {
			logrus.Infof("excluded ")
		}
		return false
	}
	if !words.Contains(word, cmd.includes, verbose) {
		if verbose {
			logrus.Infof("contains ")
		}
		return false
	}
	if !words.ContainsExact(word, cmd.exactPosition, verbose) {
		if verbose {
			logrus.Infof("exact ")
		}
		return false
	}
	if verbose {
		logrus.Infof("word %s should be included ", word.Text)
	}
	return true
}

// func (foo *rankCommand) letterWeights(words []string) map[rune]int {
// 	tuples := runeCounts(words)
// 	for i := range tuples {
// 		tuples[i].count = (26 - i) / 5
// 		if !cmd.vowels && strings.ContainsRune(vowels, tuples[i].letter) {
// 			tuples[i].count = 1
// 		}
// 	}
// 	scores := make(map[rune]int)
// 	for _, t := range tuples {
// 		scores[t.letter] = t.count
// 	}
// 	return scores
// }

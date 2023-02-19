package command

import (
	"fmt"
	"os"
	"strings"

	"wordle/words"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// const (
// 	vowels = "aeiou"
// )

type rankCommand struct {
	dictionary    string
	excluded      string
	matched       string
	wrongPosition []string
	exactPosition string
	targets       []string
	// vowels        bool
}

func NewRank(app Commander) (string, RunFunc) {
	var args rankCommand
	cmd := app.Command("rank", "rank words based on strategy")
	cmd.Flag("dictionary", "file from which to read dictionary words").Default("five-letter-words.txt").StringVar(&args.dictionary)
	cmd.Flag("match", "show wods that contain all of these letters").Short('m').StringVar(&args.matched)
	cmd.Flag("include", "letters that must be present, but not in position specified").Short('i').StringsVar(&args.wrongPosition)
	cmd.Flag("exclude", "letters that are not part of the target word").Short('e').StringVar(&args.excluded)
	cmd.Flag("exact", "word containing letters in exact positions using space for unknown. Example ' a  e'").Short('x').StringVar(&args.exactPosition)
	cmd.Flag("target", "words for which verbose logging should be enabled").Short('t').StringsVar(&args.targets)
	// cmd.Flag("vowels", "rank vowels more important").BoolVar(&args.vowels)
	return cmd.FullCommand(), args.Run
}

func (cmd *rankCommand) Run() error {
	// TODO: configure log level
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel)

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
	if cmd.matched != "" && !words.ValidExclude.MatchString(cmd.matched) {
		return fmt.Errorf("match string [%s] contains invalid character(s)", cmd.matched)
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
	for _, word := range cmd.wrongPosition {
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

func updateLogger(word *words.Word, targets []string) {
	for _, w := range targets {
		if word.Text == w {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			return
		}
	}
	// TODO: restore original level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func (cmd *rankCommand) filterWord(word *words.Word) bool {
	updateLogger(word, cmd.targets)
	if !words.Matched(word, cmd.matched) {
		log.Debug().Str("word", word.Text).Msg("match filter")
		return false
	}
	if words.Excluded(word, cmd.excluded) {
		log.Debug().Str("word", word.Text).Msg("excluded filter")
		return false
	}
	if !words.Contains(word, cmd.wrongPosition) {
		log.Debug().Str("word", word.Text).Msg("contains filter")
		return false
	}
	if !words.ContainsExact(word, cmd.exactPosition) {
		log.Debug().Str("word", word.Text).Msg("exact filter")
		return false
	}
	log.Debug().Str("word", word.Text).Msg("included")
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

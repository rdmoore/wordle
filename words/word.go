package words

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"wordle/file"

	"github.com/sirupsen/logrus"
)

var (
	ValidWord    = regexp.MustCompile(`^[a-z]{5}$`)
	ValidExclude = regexp.MustCompile(`^[a-z]+$`)
	ValidInclude = regexp.MustCompile(`^[. a-z]{5}$`)
)

type Word struct {
	Text    string
	AsRunes []rune
}

func (w *Word) String() string {
	return w.Text
}

func Scores(words []*Word, counts RuneCounts) []Score {
	scores := make([]Score, 0, len(words))
	for _, w := range words {
		found := make(map[rune]struct{})
		positionScore := make([]int, len(w.Text))
		for i, ch := range w.Text {
			if i == 4 && ch == 's' {
				continue
			}
			if _, ok := found[ch]; ok {
				continue
			}
			found[ch] = struct{}{}
			positionScore[i] = counts[ch][i]
		}
		total := 0
		for _, s := range positionScore {
			total += s
		}

		scores = append(scores, Score{
			Word:       w.Text,
			ByPosition: positionScore,
			Total:      total,
		})
	}
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Total == scores[j].Total {
			return scores[i].Word < scores[j].Word
		}
		return scores[i].Total > scores[j].Total
	})
	return scores
}

func New(text string) (*Word, error) {
	text = strings.TrimSpace(text)
	if !ValidWord.MatchString(text) {
		return nil, fmt.Errorf("invalid dictionary word [%s]", text)
	}
	return &Word{
		Text:    text,
		AsRunes: []rune(text),
	}, nil
}

func LoadWords(filename string, filter func(*Word) bool) ([]*Word, error) {
	var words []*Word
	err := file.ForEachLine(filename, func(text string) error {
		word, err := New(text)
		if err != nil {
			return err
		}
		if filter(word) {
			words = append(words, word)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load '%s': %w", filename, err)
	}
	return words, err
}

func ContainsExact(word *Word, required string, verbose bool) bool {
	for i, ch := range required {
		switch ch {
		case ' ', '.':
			continue // supported 'nil' characters
		}
		if word.AsRunes[i] != ch {
			if verbose {
				logrus.Debugf("word [%s] does not contain required char [%c] at position %d", word, ch, i)
			}
			return false
		}
	}
	return true
}

// Contains returns true if the word contains characters listed in the required strings
// as long as the characters are not in the explicit position of the character in one
// of the required. Space (' ') and period ('.') in the "required" strings are ignored
// to allow the caller to supply positional 'nil' characters.
func Contains(word *Word, required []string, verbose bool) bool {
	for _, req := range required {
		for i, r := range req {
			switch r {
			case ' ', '.':
				continue // supported 'nil' characters
			}
			if !strings.ContainsRune(word.Text, r) {
				if verbose {
					logrus.Debugf("word [%s] does not contain required characters [%c]", word, r)
				}
				return false
			}
			if word.AsRunes[i] == r {
				if verbose {
					logrus.Debugf("word [%s] contains [%c] in position [%d]", word, r, i)
				}
				return false
			}
		}
	}
	return true
}

func Excluded(word *Word, excluded string, verbose bool) bool {
	if excluded == "" || !strings.ContainsAny(word.Text, excluded) {
		return false
	}
	if verbose {
		logrus.Debugf("word [%s] contains one or more excluded characters [%s]", word, excluded)
	}
	return true
}

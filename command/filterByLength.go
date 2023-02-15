package command

import (
	"fmt"
	"wordle/file"
)

type filterCommand struct {
	input  string
	output string
	length int
}

func NewFilter(app Commander) (string, RunFunc) {
	var args filterCommand
	cmd := app.Command("filter", "filter five letter words")
	cmd.Flag("input", "file from which to read words").Default("words_alpha.txt").StringVar(&args.input)
	cmd.Flag("output", "file to write results").Default("five-letter-words.txt").StringVar(&args.output)
	cmd.Flag("len", "length of words to select").Default("5").IntVar(&args.length)
	return cmd.FullCommand(), args.Run
}

func (cmd *filterCommand) Run() error {
	fmt.Printf("hello wordle!\n")
	var matched []string
	err := file.ForEachLine(cmd.input, func(text string) error {
		if len(text) == cmd.length {
			matched = append(matched, text)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read words: %w", err)
	}
	err = file.SaveLines(cmd.output, matched)
	if err != nil {
		return fmt.Errorf("failed to write words: %w", err)
	}
	fmt.Printf("found %d five letter words", len(matched))
	return nil
}

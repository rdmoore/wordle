package main

import (
	"log"
	"os"

	"wordle/command"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("wordle", "wordle fun")
	_ = app.HelpFlag.Short('h')
	commands := make(command.Set)
	commands.Add(command.NewFilter(app))
	commands.Add(command.NewRank(app))

	cmd, ok := commands[kingpin.MustParse(app.Parse(os.Args[1:]))]
	if !ok {
		panic("command not found in map")
	}
	if err := cmd(); err != nil {
		log.Fatalf("failed: %s", err.Error())
	}
}

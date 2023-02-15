package command

import "gopkg.in/alecthomas/kingpin.v2"

type Commander interface {
	Command(name, help string) *kingpin.CmdClause
}

type RunFunc func() error

type (
	Set map[string]RunFunc
)

func (s Set) Add(command string, f RunFunc) {
	s[command] = f
}

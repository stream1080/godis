package database

import "strings"

var cmdTable = make(map[string]*Command)

type Command struct {
	executor ExecFunc
	arity    int
}

func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &Command{
		executor: executor,
		arity:    arity,
	}
}

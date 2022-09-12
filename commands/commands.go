package commands

import (
	"errors"
	aw "github.com/deanishe/awgo"
)

type command interface {
	exec() error
}

func Handle(wf *aw.Workflow, cmd string) error {
	var commands = map[string]command{
		"currency":  &currency{wf, nil},
		"hash":      &hash{wf},
		"timestamp": &timestamp{wf},
	}

	if c, ok := commands[cmd]; ok {
		return c.exec()
	}
	return errors.New("command not found")
}

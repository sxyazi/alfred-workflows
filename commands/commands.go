package commands

import (
	"errors"
	"fmt"
	aw "github.com/deanishe/awgo"
	"strings"
)

type command interface {
	general([]string) error
	universal(string, []string) (string, error)
}

func Handle(wf *aw.Workflow, args []string) error {
	cmd := args[0]
	var commands = map[string]command{
		"currency":  &currency{wf, nil},
		"datetime":  &datetime{wf},
		"decode":    &decode{wf},
		"encode":    &encode{wf},
		"hash":      &hash{wf},
		"timestamp": &timestamp{wf},
		"uuid":      &uuid{wf},
	}

	var sub = ""
	if s := strings.SplitN(cmd, ":", 2); len(s) == 2 {
		cmd, sub = s[0], s[1]
	}

	if c, ok := commands[cmd]; !ok {
		return errors.New("command not found: " + cmd)
	} else if sub == "" {
		return c.general(args[1:])
	} else if op, err := c.universal(sub, args[1:]); err != nil {
		return err
	} else {
		fmt.Print(op)
		return nil
	}
}

package main

import (
	"github.com/deanishe/awgo"
	"github.com/sxyazi/alfred-workflows/commands"
)

var wf *aw.Workflow

func run() {
	if err := commands.Handle(wf, wf.Args()); err != nil {
		panic(err)
	}
}

func main() {
	wf = aw.New()
	wf.Run(run)
}

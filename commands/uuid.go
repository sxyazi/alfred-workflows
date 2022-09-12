package commands

import (
	"errors"
	aw "github.com/deanishe/awgo"
	"sync"
)
import uuid_ "github.com/google/uuid"

type uuid struct {
	wf *aw.Workflow
}

func (u *uuid) general(args []string) error {
	wg := sync.WaitGroup{}
	wg.Add(8)

	for i := 0; i < 8; i++ {
		go func() {
			defer wg.Done()

			s := uuid_.New().String()
			u.wf.
				NewItem(s).
				Arg(s).
				Valid(true).
				Icon(&aw.Icon{Value: u.wf.Dir() + "/static/code.png"})
		}()
	}

	wg.Wait()
	u.wf.SendFeedback()
	return nil
}

func (u *uuid) universal(act string, args []string) (string, error) {
	return "", errors.New("unsupported action")
}

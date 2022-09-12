package commands

import (
	"encoding/base64"
	"errors"
	aw "github.com/deanishe/awgo"
	"sync"
)

type encode struct {
	wf *aw.Workflow
}

func (e *encode) general(args []string) error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	var b64 string
	go func() {
		defer wg.Done()
		b64 = base64.StdEncoding.EncodeToString([]byte(args[0]))
	}()

	wg.Wait()
	e.wf.
		NewItem("Base64").
		Subtitle(b64).
		Arg(b64).
		Valid(true).
		Icon(&aw.Icon{Value: e.wf.Dir() + "/static/base64.png"})

	e.wf.SendFeedback()
	return nil
}

func (e *encode) universal(act string, args []string) (string, error) {
	switch act {
	case "base64":
		return base64.StdEncoding.EncodeToString([]byte(args[0])), nil
	}
	return "", errors.New("unsupported action")
}

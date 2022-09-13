package commands

import (
	"encoding/base64"
	"errors"
	aw "github.com/deanishe/awgo"
	"github.com/sxyazi/alfred-workflows/utils"
	"net/url"
	"sync"
)

type decode struct {
	wf *aw.Workflow
}

func (d *decode) general(args []string) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var query, b64 string
	go func() {
		defer wg.Done()
		b64 = utils.Value(url.QueryUnescape(args[0]))
	}()
	go func() {
		defer wg.Done()
		b64 = string(utils.Value(base64.StdEncoding.DecodeString(args[0])))
	}()

	wg.Wait()
	if query != "" {
		d.wf.
			NewItem("URL").
			Subtitle(query).
			Arg(query).
			Valid(true).
			Icon(&aw.Icon{Value: d.wf.Dir() + "/static/link.png"})
	}
	if b64 != "" {
		d.wf.
			NewItem("Base64").
			Subtitle(b64).
			Arg(b64).
			Valid(true).
			Icon(&aw.Icon{Value: d.wf.Dir() + "/static/base64.png"})
	}

	d.wf.SendFeedback()
	return nil
}

func (d *decode) universal(act string, args []string) (string, error) {
	switch act {
	case "url":
		return url.QueryUnescape(args[0])
	case "base64":
		return string(utils.Value(base64.StdEncoding.DecodeString(args[0]))), nil
	}
	return "", errors.New("unsupported action")
}

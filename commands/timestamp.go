package commands

import (
	aw "github.com/deanishe/awgo"
	"strconv"
	"time"
)

type timestamp struct {
	wf *aw.Workflow
}

func (t *timestamp) exec() error {
	ts := time.Now().UnixMilli()

	t.wf.
		NewItem("Timestamp (s)").
		Subtitle(strconv.FormatInt(ts/1000, 10)).
		Arg(strconv.FormatInt(ts/1000, 10)).
		Valid(true).
		Icon(aw.IconNote)

	t.wf.
		NewItem("Timestamp (ms)").
		Subtitle(strconv.FormatInt(ts, 10)).
		Arg(strconv.FormatInt(ts, 10)).
		Valid(true).
		Icon(aw.IconNote)

	t.wf.SendFeedback()
	return nil
}

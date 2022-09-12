package commands

import (
	"errors"
	aw "github.com/deanishe/awgo"
	"github.com/sxyazi/alfred-workflows/utils"
	"strconv"
	"strings"
	"time"
)

type timestamp struct {
	wf *aw.Workflow
}

func (t *timestamp) general(args []string) error {
	var arg string
	if len(args) >= 1 {
		arg = strings.TrimSpace(args[0])
	}

	var r *time.Time
	if arg == "" {
		r = utils.Ptr(time.Now())
	} else if r = utils.ParseTime(arg); r == nil {
		return errors.New("invalid date")
	}

	icon := &aw.Icon{Value: t.wf.Dir() + "/static/clock.png"}
	s := strconv.FormatInt(r.Unix(), 10)
	ms := strconv.FormatInt(r.UnixMilli(), 10)
	t.wf.
		NewItem("Timestamp (s)").
		Subtitle(s).
		Arg(s).
		Valid(true).
		Icon(icon)

	t.wf.
		NewItem("Timestamp (ms)").
		Subtitle(ms).
		Arg(ms).
		Valid(true).
		Icon(icon)

	t.wf.SendFeedback()
	return nil
}

func (t *timestamp) universal(act string, args []string) (string, error) {
	return "", errors.New("unsupported action")
}

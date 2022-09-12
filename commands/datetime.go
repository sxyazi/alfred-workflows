package commands

import (
	"errors"
	"fmt"
	aw "github.com/deanishe/awgo"
	"github.com/sxyazi/alfred-workflows/utils"
	"strings"
	"time"
)

type datetime struct {
	wf *aw.Workflow
}

func (d *datetime) general(args []string) error {
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

	icon := &aw.Icon{Value: fmt.Sprintf("%s/static/calendar/%d.png", d.wf.Dir(), r.Day())}
	d.wf.
		NewItem(r.Format("2006-01-02 15:04:05")).
		Arg(r.Format("2006-01-02 15:04:05")).
		Valid(true).
		Icon(icon)

	d.wf.
		NewItem(r.Format("2006-01-02")).
		Arg(r.Format("2006-01-02")).
		Valid(true).
		Icon(icon)

	d.wf.
		NewItem(r.Format("15:04:05")).
		Arg(r.Format("15:04:05")).
		Valid(true).
		Icon(icon)

	d.wf.
		NewItem(r.Format(time.RFC3339)).
		Arg(r.Format(time.RFC3339)).
		Valid(true).
		Icon(icon)

	d.wf.SendFeedback()
	return nil
}

func (d *datetime) universal(act string, args []string) (string, error) {
	return "", errors.New("unsupported action")
}

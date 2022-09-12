package commands

import (
	aw "github.com/deanishe/awgo"
	"github.com/sxyazi/alfred-workflows/utils"
	"sync"
)

type hash struct {
	wf *aw.Workflow
}

func (h *hash) exec() error {
	s := h.wf.Args()[1]

	wg := sync.WaitGroup{}
	wg.Add(4)

	var md5, sha1, sha256, sha512 string
	go func() {
		defer wg.Done()
		md5 = utils.StrMd5(s)
	}()
	go func() {
		defer wg.Done()
		sha1 = utils.StrSha1(s)
	}()
	go func() {
		defer wg.Done()
		sha256 = utils.StrSha256(s)
	}()
	go func() {
		defer wg.Done()
		sha512 = utils.StrSha512(s)
	}()

	wg.Wait()
	h.wf.
		NewItem("MD5").
		Subtitle(md5).
		Arg(md5).
		Valid(true).
		Icon(aw.IconNote)

	h.wf.
		NewItem("SHA1").
		Subtitle(sha1).
		Arg(sha1).
		Valid(true).
		Icon(aw.IconNote)

	h.wf.
		NewItem("SHA256").
		Subtitle(sha256).
		Arg(sha256).
		Valid(true).
		Icon(aw.IconNote)

	h.wf.
		NewItem("SHA512").
		Subtitle(sha512).
		Arg(sha512).
		Valid(true).
		Icon(aw.IconNote)

	h.wf.SendFeedback()
	return nil
}

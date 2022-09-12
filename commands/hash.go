package commands

import (
	"errors"
	aw "github.com/deanishe/awgo"
	"github.com/sxyazi/alfred-workflows/utils"
	"sync"
)

type hash struct {
	wf *aw.Workflow
}

func (h *hash) general(args []string) error {
	wg := sync.WaitGroup{}
	wg.Add(4)

	var md5, sha1, sha256, sha512 string
	go func() {
		defer wg.Done()
		md5 = utils.StrMd5(args[0])
	}()
	go func() {
		defer wg.Done()
		sha1 = utils.StrSha1(args[0])
	}()
	go func() {
		defer wg.Done()
		sha256 = utils.StrSha256(args[0])
	}()
	go func() {
		defer wg.Done()
		sha512 = utils.StrSha512(args[0])
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

func (h *hash) universal(act string, args []string) (string, error) {
	switch act {
	case "md5":
		return utils.StrMd5(args[0]), nil
	case "sha1":
		return utils.StrSha1(args[0]), nil
	case "sha256":
		return utils.StrSha256(args[0]), nil
	case "sha512":
		return utils.StrSha512(args[0]), nil
	}
	return "", errors.New("unsupported action")
}

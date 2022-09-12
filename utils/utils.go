package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	collect "github.com/sxyazi/go-collection"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Ptr[T any](v T) *T {
	return &v
}

func Value[T any](first T, args ...any) T {
	return first
}

func StrMd5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func StrSha1(s string) string {
	hash := sha1.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func StrSha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func StrSha512(s string) string {
	hash := sha512.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func NewClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 200
	t.MaxConnsPerHost = 10
	t.MaxIdleConnsPerHost = 10
	t.IdleConnTimeout = 5 * time.Minute

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: t,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

var Client = NewClient()

func DeepValue[V any](v any, key string) (V, error) {
	var err error
	var zero V
	for _, k := range strings.Split(key, ".") {
		v, err = collect.AnyGet[any](v, k)
		if err != nil {
			return zero, err
		}
	}

	if v, ok := v.(V); ok {
		return v, nil
	}

	switch any(zero).(type) {
	case float64:
		if n, err := strconv.ParseFloat(fmt.Sprintf("%f", v), 64); err == nil {
			return any(n).(V), nil
		}
	case int:
		if n, err := strconv.ParseFloat(fmt.Sprintf("%f", v), 64); err == nil {
			return any(int(n)).(V), nil
		}
	}
	return zero, errors.New("type mismatch")
}

func HttpPost(u string, d io.Reader) ([]byte, error) {
	req, err := http.NewRequest("POST", u, d)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1321.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, err
}

var dateFormats = []string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,

	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04",
	"2006/01/02 15:04:05",

	"06-1-2",
	"2006-1-2",
	"2006-01-02",

	"06/1/2",
	"2006/1/2",
	"2006/01/02",

	"15:4:5",
	"15:04:05",
}

func ParseTime(s string) *time.Time {
	if n, err := strconv.ParseInt(s, 10, 64); err == nil && n > 0 {
		if len(s) == 10 {
			return Ptr(time.Unix(n, 0))
		} else if len(s) == 13 {
			return Ptr(time.UnixMilli(n))
		}
	}

	for _, f := range dateFormats {
		if t, err := time.Parse(f, s); err == nil {
			return &t
		}
	}
	return nil
}

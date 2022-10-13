package limits

import (
	"context"
	"github.com/BOOMfinity/go-utils/rate"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var paths = []string{"channels", "guilds"}

type Manager struct {
	buckets map[string]*rate.Limiter
	prefix  string
	m       sync.RWMutex
}

// ParseBucketPath
//
// Reference: https://github.com/diamondburned/arikawa/blob/v3/api/rate/majors.go
func (v *Manager) ParseBucketPath(url string) string {
	url = strings.TrimPrefix(url, v.prefix)
	url = strings.SplitN(url, "?", 2)[0]

	parts := strings.Split(url, "/")
	if len(parts) < 1 {
		return url
	}

	parts = parts[1:]

	var skip int

	for _, part := range paths {
		if part == parts[0] {
			skip = 2
			break
		}
	}
	skip++

	for ; skip < len(parts); skip += 2 {
		if _, err := strconv.Atoi(parts[skip]); err == nil {
			parts[skip] = ""
			continue
		}
		if StringIsEmojiOnly(parts[skip]) {
			parts[skip] = ""
			continue
		}
		if StringIsCustomEmoji(parts[skip]) {
			parts[skip] = ""
			continue
		}
	}

	url = strings.Join(parts, "/")
	return url
}

func (v *Manager) ParseRequest(req *fasthttp.Request, res *fasthttp.Response) error {
	limitRaw := string(res.Header.Peek("x-ratelimit-limit"))
	remainingRaw := string(res.Header.Peek("x-ratelimit-remaining"))
	resetAfterRaw := string(res.Header.Peek("x-ratelimit-reset-after"))
	id := string(res.Header.Peek("x-ratelimit-bucket"))
	if id == "" || limitRaw == "" || remainingRaw == "" || resetAfterRaw == "" {
		return nil
	}
	//println(resetAfterRaw, id, limitRaw, remainingRaw)
	limit, err := strconv.ParseInt(limitRaw, 10, 64)
	remaining, err := strconv.ParseInt(remainingRaw, 10, 64)
	resetAfter, err := strconv.ParseFloat(resetAfterRaw, 64)
	if err != nil {
		return err
	}
	if remaining != (limit - 1) {
		return nil
	}
	path := v.ParseBucketPath(string(req.URI().FullURI()))
	bucket := v.get(path)
	if bucket == nil {
		bucket = rate.NewLimiterInit(time.Duration(resetAfter)*time.Second, int(limit), 1)
		v.m.Lock()
		v.buckets[path] = bucket
		v.m.Unlock()
	} else {
		bucket.Update(time.Duration(resetAfter)*time.Second, int(limit))
	}
	return nil
}

func (v *Manager) get(path string) *rate.Limiter {
	v.m.RLock()
	defer v.m.RUnlock()
	return v.buckets[path]
}

func (v *Manager) Wait(ctx context.Context, url string) error {
	v.m.RLock()
	defer v.m.RUnlock()
	bucket := v.get(v.ParseBucketPath(url))
	if bucket == nil {
		return nil
	}
	return bucket.Wait(ctx)
}

func NewManager(prefix string) *Manager {
	return &Manager{
		buckets: map[string]*rate.Limiter{},
		prefix:  prefix,
	}
}

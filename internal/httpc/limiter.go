package httpc

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BOOMfinity/go-utils/ubytes"
	"github.com/BOOMfinity/golog/v2"
	"github.com/andersfylling/snowflake/v5"
	"github.com/forPelevin/gomoji"
	"github.com/valyala/fasthttp"

	"github.com/BOOMfinity/bfcord"
	"github.com/BOOMfinity/bfcord/utils"
)

func parseBucket(url string) (bucket string) {
	url = strings.TrimPrefix(url, bfcord.APIUrl+"/"+bfcord.APIVersion+"/")
	segments := strings.Split(url, "/")

	skip := 0

	if segments[skip] == "guilds" ||
		segments[skip] == "channels" ||
		segments[skip] == "webhooks" {
		skip = 2
	}

	for ; skip < len(segments); skip += 1 {

		if _, err := snowflake.ParseSnowflakeUint(segments[skip], 10); err == nil {
			segments[skip] = ""
		}

		if strings.Contains(segments[skip], ":") {
			segments[skip] = ""
		}

		if stringIsEmoji(segments[skip]) {
			segments[skip] = ""
		}

	}

	return bucket + path.Join(segments...)
}

func stringIsEmoji(str string) bool {
	runes := []rune(str)
	if len(runes) == 1 ||
		len(runes) == 2 {
		if gomoji.ContainsEmoji(str) {
			return true
		}
	}

	return false
}

type bucket struct {
	name string

	remaining uint64
	reset     time.Time
	lock      utils.CustomMutex
}

type limiter struct {
	log     golog.Logger
	buckets map[string]*bucket
	mapLock sync.Mutex

	global time.Time
}

func (l *limiter) getBucket(url string) *bucket {
	name := parseBucket(url)
	l.mapLock.Lock()
	b := l.buckets[name]
	if b == nil {
		b = &bucket{
			name: name,
			lock: utils.NewMutex(),
		}
		l.buckets[name] = b
	}
	l.mapLock.Unlock()
	return b
}

func (l *limiter) Acquire(ctx context.Context, url string) error {
	bucket := l.getBucket(url)
	bucket.lock.Lock(ctx)
	defer bucket.lock.Unlock()

	now := time.Now()
	until := time.Time{}

	if bucket.remaining == 0 && now.Before(bucket.reset) {
		until = bucket.reset
	} else {
		until = l.global
	}

	if now.Before(until) {
		l.log.Trace().Param("bucket", bucket.name).Send("Have to wait %s before allowing to acquire the bucket", until.Sub(now).String())
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(until.Sub(now)):
		}
	}

	if bucket.remaining > 0 {
		bucket.remaining--
	}

	return nil
}

func (l *limiter) Release(url string, headers *fasthttp.ResponseHeader) error {
	bucket := l.getBucket(url)
	if bucket == nil {
		return nil
	}
	defer bucket.lock.TryUnlock()

	var (
		global = len(headers.Peek("X-RateLimit-Global")) > 0

		remaining  = ubytes.ToString(headers.Peek("X-RateLimit-Remaining"))
		reset      = ubytes.ToString(headers.Peek("X-RateLimit-Reset"))
		retryAfter = ubytes.ToString(headers.Peek("Retry-After"))
	)

	switch {
	case retryAfter != "":
		i, err := strconv.Atoi(retryAfter)
		if err != nil {
			return fmt.Errorf("failed to convert string to int (retry-after): %w", err)
		}

		at := time.Now().Add(time.Duration(i) * time.Second)

		if global {
			l.global = at
		} else {
			bucket.reset = at
		}
	case reset != "":
		unix, err := strconv.ParseFloat(reset, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float64 from X-RateLimit-Reset header (%s): %w", reset, err)
		}

		secs := int64(unix)
		nsecs := int64((unix - float64(secs)) * float64(time.Second))

		// added extra delay because Discord is stupid
		bucket.reset = time.Unix(secs, nsecs).Add(250 * time.Millisecond)
	case remaining != "":
		u, err := strconv.ParseUint(remaining, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint64 from X-RateLimit-Remaining header (%s): %w", remaining, err)
		}

		bucket.remaining = u
	}
	return nil
}

package httpc

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/BOOMfinity/bfcord"
	"github.com/BOOMfinity/go-utils/rate"
	"github.com/BOOMfinity/golog/v2"
)

func ResolvePath(segments ...string) string {
	return bfcord.APIUrl + "/" + bfcord.APIVersion + "/" + strings.Join(segments, "/")
}

type Client struct {
	token   string
	log     golog.Logger
	limiter *rate.Limiter
	buckets *limiter
	id      atomic.Uint64
}

func (c *Client) CustomGlobalLimiter(requests int) {
	c.limiter = rate.NewLimiter(1*time.Second, requests)
}

func NewClient(token string, log golog.Logger) *Client {
	if log == nil {
		log = golog.New("api")
	}
	return &Client{
		token: token,
		log:   log,
		buckets: &limiter{
			buckets: map[string]*bucket{},
			log:     log.Module("buckets"),
		},
		limiter: rate.NewLimiter(1*time.Second, 50),
	}
}

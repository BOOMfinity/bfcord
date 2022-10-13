package api

import (
	"github.com/BOOMfinity/golog"
	"time"
)

type RequestData struct {
	logger           golog.Logger
	prefix           string
	authHeaderPrefix string
	timeout          time.Duration
	retryDelay       time.Duration
	retries          uint8
}

type Option func(v *RequestData)

func WithLogger(log golog.Logger) Option {
	return func(v *RequestData) {
		v.logger = log
	}
}

func WithAuthHeaderPrefix(prefix string) Option {
	return func(v *RequestData) {
		v.authHeaderPrefix = prefix
	}
}

func WithPrefix(prefix string) Option {
	return func(v *RequestData) {
		v.prefix = prefix
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(v *RequestData) {
		v.timeout = timeout
	}
}

func WithRetries(retries uint8) Option {
	return func(v *RequestData) {
		v.retries = retries
	}
}

func WithRetryDelay(delay time.Duration) Option {
	return func(v *RequestData) {
		v.retryDelay = delay
	}
}

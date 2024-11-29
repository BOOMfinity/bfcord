package utils

import (
	"context"
	"time"
)

func NewMutex() CustomMutex {
	return &customMutexImpl{
		ch: make(chan bool, 1),
	}
}

type CustomMutex interface {
	Lock(ctx context.Context) error
	Unlock()
	TryLock(ctx context.Context) bool
	TryUnlock() bool
}

type customMutexImpl struct {
	ch chan bool
}

func (m *customMutexImpl) Lock(ctx context.Context) error {
	select {
	case m.ch <- true:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *customMutexImpl) Unlock() {
	select {
	case <-m.ch:
		return
	default:
		panic("trying to unlock the already unlocked mutex!")
	}
}

func (m *customMutexImpl) TryLock(ctx context.Context) bool {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, 25*time.Millisecond)
		defer cancel()
	}
	select {
	case m.ch <- true:
		return true
	case <-ctx.Done():
		return false
	}
}

func (m *customMutexImpl) TryUnlock() bool {
	select {
	case <-m.ch:
		return true
	default:
		return false
	}
}

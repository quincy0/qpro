package rescue

import (
	"context"
	"github.com/quincy0/qpro/qLog"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	qLog.MustNew()
}

func TestRecover(t *testing.T) {
	var count int32
	assert.NotPanics(t, func() {
		defer Recover(func() {
			atomic.AddInt32(&count, 2)
		}, func() {
			atomic.AddInt32(&count, 3)
		})

		panic("hello")
	})

	assert.Equal(t, int32(5), atomic.LoadInt32(&count))
}

func TestRecoverCtx(t *testing.T) {
	var count int32
	assert.NotPanics(t, func() {
		defer RecoverCtx(context.Background(), func() {
			atomic.AddInt32(&count, 2)
		}, func() {
			atomic.AddInt32(&count, 3)
		})

		panic("hello")
	})

	assert.Equal(t, int32(5), atomic.LoadInt32(&count))
}

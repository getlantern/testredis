package testredis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {

	r, err := OpenUnstarted()
	assert.NoError(t, err)

	r.Start()

	cl := r.Client()

	ctx := context.Background()
	pong := cl.Ping(context.Background())
	val, err := pong.Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", val)

	cl.Set(ctx, "key", "val", 10000*time.Second)
	cmd := cl.Get(ctx, "key")

	vals, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(vals))
}

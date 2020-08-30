package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestIDOk(t *testing.T) {
	expected := "request_id"

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxRequestID, expected)

	actual := GetRequestID(ctx)
	assert.Equal(t, expected, actual)
}

func TestGetRequestIDEmpty(t *testing.T) {
	ctx := context.Background()

	actual := GetRequestID(ctx)
	assert.Equal(t, "", actual)
}

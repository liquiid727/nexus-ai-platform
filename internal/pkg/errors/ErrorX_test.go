package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewAndErrorString(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		reason  string
		format  string
		args    []any
		wantMsg string
	}{
		{"basic", 400, "Bad", "bad: %s", []any{"x"}, "bad: x"},
		{"empty-msg", 200, "OK", "", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.code, tt.reason, tt.format, tt.args...)
			assert.Equal(t, tt.code, e.Code)
			assert.Equal(t, tt.reason, e.Reason)
			assert.Equal(t, tt.wantMsg, e.Message)
			out := e.Error()
			assert.Contains(t, out, "code = ")
			assert.Contains(t, out, "reason = ")
			assert.Contains(t, out, "message = ")
		})
	}
}

func TestWithMessageAndWithMetadata(t *testing.T) {
	e := New(500, "Internal", "init")

	e.WithMessage("hello %s", "world")
	assert.Equal(t, "hello world", e.Message)

	md := map[string]string{"a": "1", "b": "2"}
	e.WithMetadata(md)
	assert.Equal(t, md, e.Metadata)
}

func TestKVEdgeCases(t *testing.T) {
	e := New(500, "Internal", "")

	e.KV("k1", "v1", "k2") // odd count, last ignored
	assert.Equal(t, "v1", e.Metadata["k1"])
	_, ok := e.Metadata["k2"]
	assert.False(t, ok, "odd kv should ignore last key without value")

	e.KV() // no-op
	assert.Equal(t, 1, len(e.Metadata))
}

func TestWithRequestID(t *testing.T) {
	e := New(200, "OK", "")
	id := "req-123"
	e.WithRequestID(id)
	assert.Equal(t, id, e.Metadata["X-nexus-gateway-request_id"])
}

func TestIs(t *testing.T) {
	a := New(404, "NotFound", "x")
	b := New(404, "NotFound", "y")
	c := New(400, "BadRequest", "z")

	assert.True(t, a.Is(b))
	assert.False(t, a.Is(c))
	assert.False(t, a.Is(errors.New("raw")))
}

func TestCodeAndReason(t *testing.T) {
	assert.Equal(t, 200, Code(nil))
	assert.Equal(t, " ", Reason(nil))

	e := New(401, "Unauthorized", "nope")
	assert.Equal(t, 401, Code(e))
	assert.Equal(t, "Unauthorized", Reason(e))

	raw := errors.New("boom")
	assert.Equal(t, ErrInternal.Code, Code(raw))
	assert.Equal(t, ErrInternal.Reason, Reason(raw))
}

func TestFromErrorVariants(t *testing.T) {
	assert.Nil(t, FromError(nil))

	e := New(403, "Forbidden", "stop")
	got := FromError(e)
	assert.Equal(t, e, got)

	raw := errors.New("x")
	got2 := FromError(raw)
	assert.Equal(t, ErrInternal.Code, got2.Code)
	assert.Equal(t, ErrInternal.Reason, got2.Reason)
	assert.Contains(t, got2.Message, "x")

	st := status.New(codes.NotFound, "missing")
	info := &errdetails.ErrorInfo{
		Reason:   "ResourceMissing",
		Domain:   "gateway",
		Metadata: map[string]string{"rid": "1"},
	}
	stWith, err := st.WithDetails(info)
	assert.NoError(t, err)
	grpcErr := stWith.Err()

	got3 := FromError(grpcErr)
	assert.Equal(t, int(codes.NotFound), got3.Code)
	// Message comes from err.Error() per implementation
	assert.Contains(t, got3.Message, "NotFound")
	assert.Contains(t, got3.Message, "missing")
	assert.Equal(t, "ResourceMissing", got3.Reason)
	assert.Equal(t, "1", got3.Metadata["rid"])
}

func TestGRPCStatus(t *testing.T) {
	e := New(500, "Internal", "")
	assert.Nil(t, e.GRPCStatus())
}

func TestConcurrencyKV(t *testing.T) {
	// Concurrency across instances should be safe
	N := 50
	ch := make(chan *ErrorX, N)
	for i := 0; i < N; i++ {
		go func(i int) {
			e := New(200, "OK", "")
			e.KV("k"+string(rune(i)), "v"+string(rune(i)))
			ch <- e
		}(i)
	}
	for i := 0; i < N; i++ {
		select {
		case e := <-ch:
			assert.Equal(t, 1, len(e.Metadata))
		}
	}
}

func BenchmarkKV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := New(200, "OK", "")
		for j := 0; j < 1000; j++ {
			e.KV("k", "v")
		}
	}
}

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/internal"

	"github.com/stretchr/testify/assert"
)

func TestContextWithSpan(t *testing.T) {
	want, ctx := StartSpanFromContext(context.TODO(), "http.request")
	ctx = ContextWithSpan(context.TODO(), want)
	got, ok := ctx.Value(internal.ActiveSpanKey).(ddtrace.Span)
	assert := assert.New(t)
	assert.True(ok)
	assert.Equal(got, want)
}

func TestStartSpanFromContextRace(t *testing.T) {
	_, stop := startTestTracer(t)
	defer stop()

	// Start 100 goroutines that create child spans with StartSpanFromContext in parallel,
	// with a shared options slice. The child spans should get parented to the correct spans
	const contextKey = "key"
	const numContexts = 100
	options := make([]StartSpanOption, 0, 3)
	outputValues := make(chan uint64, numContexts)
	var expectedTraceIDs []uint64
	for i := 0; i < numContexts; i++ {
		parent, childCtx := StartSpanFromContext(context.Background(), "parent")
		expectedTraceIDs = append(expectedTraceIDs, parent.Context().TraceID())
		go func() {
			span, _ := StartSpanFromContext(childCtx, "testoperation", options...)
			defer span.Finish()
			outputValues <- span.Context().TraceID()
		}()
		parent.Finish()
	}

	// collect the outputs
	var outputs []uint64
	for i := 0; i < numContexts; i++ {
		outputs = append(outputs, <-outputValues)
	}
	assert.Len(t, outputs, numContexts)
	assert.ElementsMatch(t, outputs, expectedTraceIDs)
}

func Test128(t *testing.T) {
	_, stop := startTestTracer(t)
	defer stop()

	t.Run("disable 128 bit trace ids", func(t *testing.T) {
		t.Setenv("DD_TRACE_128_BIT_TRACEID_GENERATION_ENABLED", "false")
		span, _ := StartSpanFromContext(context.Background(), "http.request")
		assert.NotZero(t, span.Context().TraceID())
		w3cCtx, ok := span.Context().(ddtrace.SpanContextW3C)
		if !ok {
			assert.Fail(t, "couldn't cast to ddtrace.SpanContextW3C")
		}
		id128 := w3cCtx.TraceID128()
		assert.Len(t, id128, 32) // ensure there are enough leading zeros
		idBytes, err := hex.DecodeString(id128)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), binary.BigEndian.Uint64(idBytes[:8])) // high 64 bits should be 0
		assert.Equal(t, span.Context().TraceID(), binary.BigEndian.Uint64(idBytes[8:]))
	})

	t.Run("enable 128 bit trace ids", func(t *testing.T) {
		// DD_TRACE_128_BIT_TRACEID_GENERATION_ENABLED is true by default
		span128, _ := StartSpanFromContext(context.Background(), "http.request")
		assert.NotZero(t, span128.Context().TraceID())
		w3cCtx, ok := span128.Context().(ddtrace.SpanContextW3C)
		if !ok {
			assert.Fail(t, "couldn't cast to ddtrace.SpanContextW3C")
		}
		id128bit := w3cCtx.TraceID128()
		assert.NotEmpty(t, id128bit)
		assert.Len(t, id128bit, 32)
		// Ensure that the lower order bits match the span's 64-bit trace id
		b, err := hex.DecodeString(id128bit)
		assert.NoError(t, err)
		assert.Equal(t, span128.Context().TraceID(), binary.BigEndian.Uint64(b[8:]))
	})
}

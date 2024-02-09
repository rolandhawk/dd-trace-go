// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	v2 "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/internal"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/log"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/telemetry"

	"github.com/stretchr/testify/assert"
)

var (
	// timeMultiplicator specifies by how long to extend waiting times.
	// It may be altered in some environments (like AppSec) where things
	// move slower and could otherwise create flaky tests.
	timeMultiplicator = time.Duration(1)

	// integration indicates if the test suite should run integration tests.
	integration bool
)

func TestMain(m *testing.M) {
	if internal.BoolEnv("DD_APPSEC_ENABLED", false) {
		// things are slower with AppSec; double wait times
		timeMultiplicator = time.Duration(2)
	}
	_, integration = os.LookupEnv("INTEGRATION")
	os.Exit(m.Run())
}

func TestTracerRuntimeMetrics(t *testing.T) {
	t.Run("on", func(t *testing.T) {
		tp := new(log.RecordLogger)
		tp.Ignore("appsec: ", telemetry.LogPrefix)
		tracer := newUnstartedTracer(WithRuntimeMetrics(), WithLogger(tp), WithDebugMode(true))
		defer tracer.Stop()
		assert.Contains(t, tp.Logs()[0], "DEBUG: Runtime metrics enabled")
	})

	t.Run("env", func(t *testing.T) {
		t.Setenv("DD_RUNTIME_METRICS_ENABLED", "true")
		tp := new(log.RecordLogger)
		tp.Ignore("appsec: ", telemetry.LogPrefix)
		tracer := newUnstartedTracer(WithLogger(tp), WithDebugMode(true))
		defer tracer.Stop()
		assert.Contains(t, tp.Logs()[0], "DEBUG: Runtime metrics enabled")
	})

	t.Run("overrideEnv", func(t *testing.T) {
		t.Setenv("DD_RUNTIME_METRICS_ENABLED", "false")
		tp := new(log.RecordLogger)
		tp.Ignore("appsec: ", telemetry.LogPrefix)
		tracer := newUnstartedTracer(WithRuntimeMetrics(), WithLogger(tp), WithDebugMode(true))
		defer tracer.Stop()
		assert.Contains(t, tp.Logs()[0], "DEBUG: Runtime metrics enabled")
	})
}

func TestTracerInjectConcurrency(t *testing.T) {
	tracer, stop := startTestTracer(t)
	defer stop()
	span, _ := StartSpanFromContext(context.Background(), "main")
	defer span.Finish()

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		i := i
		go func(val int) {
			defer wg.Done()
			span.SetBaggageItem("val", fmt.Sprintf("%d", val))

			traceContext := map[string]string{}
			_ = tracer.Inject(span.Context(), TextMapCarrier(traceContext))
		}(i)
	}

	wg.Wait()
}

// TestTracerTraceMaxSize tests a bug that was encountered in environments
// creating a large volume of spans that reached the trace cap value (traceMaxSize).
// The bug was that once the cap is reached, no more spans are pushed onto
// the buffer, yet they are part of the same trace. The trace is considered
// completed and flushed when the number of finished spans == number of spans
// in buffer. When reaching the cap, this condition might become true too
// early, and some spans in the buffer might still not be finished when flushing.
// Changing these spans at the moment of flush would (and did) cause a race
// condition.
func TestTracerTraceMaxSize(t *testing.T) {
	_, stop := startTestTracer(t)
	defer stop()

	otss, otms := traceStartSize, traceMaxSize
	traceStartSize, traceMaxSize = 3, 3
	defer func() {
		traceStartSize, traceMaxSize = otss, otms
	}()

	spans := make([]ddtrace.Span, 5)
	spans[0] = StartSpan("span0")
	spans[1] = StartSpan("span1", ChildOf(spans[0].Context()))
	spans[2] = StartSpan("span2", ChildOf(spans[0].Context()))
	spans[3] = StartSpan("span3", ChildOf(spans[0].Context()))
	spans[4] = StartSpan("span4", ChildOf(spans[0].Context()))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			spans[1].SetTag(strconv.Itoa(i), 1)
			spans[2].SetTag(strconv.Itoa(i), 1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		spans[0].Finish()
		spans[3].Finish()
		spans[4].Finish()
	}()

	wg.Wait()
}

// BenchmarkConcurrentTracing tests the performance of spawning a lot of
// goroutines where each one creates a trace with a parent and a child.
func BenchmarkConcurrentTracing(b *testing.B) {
	tracer, stop := startTestTracer(b, WithLogger(log.DiscardLogger{}), WithSampler(NewRateSampler(0)))
	defer stop()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				parent := tracer.StartSpan("pylons.request", ServiceName("pylons"), ResourceName("/"))
				defer parent.Finish()

				for i := 0; i < 10; i++ {
					tracer.StartSpan("redis.command", ChildOf(parent.Context())).Finish()
				}
			}()
		}
		wg.Wait()
	}
}

// BenchmarkTracerAddSpans tests the performance of creating and finishing a root
// span. It should include the encoding overhead.
func BenchmarkTracerAddSpans(b *testing.B) {
	tracer, stop := startTestTracer(b, WithLogger(log.DiscardLogger{}), WithSampler(NewRateSampler(0)))
	defer stop()

	for n := 0; n < b.N; n++ {
		span := tracer.StartSpan("pylons.request", ServiceName("pylons"), ResourceName("/"))
		span.Finish()
	}
}

func BenchmarkStartSpan(b *testing.B) {
	tracer, stop := startTestTracer(b, WithLogger(log.DiscardLogger{}), WithSampler(NewRateSampler(0)))
	defer stop()

	root := tracer.StartSpan("pylons.request", ServiceName("pylons"), ResourceName("/"))
	ctx := ContextWithSpan(context.TODO(), root)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s, ok := SpanFromContext(ctx)
		if !ok {
			b.Fatal("no span")
		}
		StartSpan("op", ChildOf(s.Context()))
	}
}

// startTestTracer returns a Tracer with a no-op HTTP transport
func startTestTracer(t testing.TB, opts ...StartOption) (ddtrace.Tracer, func()) {
	t.Helper()

	transport := &dummyTransport{}
	c := &http.Client{Transport: transport}
	opts = append(opts, WithLogger(log.DiscardLogger{}), WithSampler(NewRateSampler(0)), WithHTTPClient(c))
	tracer := newUnstartedTracer(opts...).(*internal.TracerV2Adapter)
	v2.SetGlobalTracer(tracer.Tracer)

	return tracer, tracer.Stop
}

type dummyTransport struct {
}

// RoundTrip implementa la interfaz RoundTripper
func (n *dummyTransport) RoundTrip(*http.Request) (*http.Response, error) {
	// notar que esta es una respuesta no inicializada
	// si se está buscando encabezados, etc.
	return &http.Response{StatusCode: http.StatusOK}, nil
}

// BenchmarkTracerStackFrames tests the performance of taking stack trace.
func BenchmarkTracerStackFrames(b *testing.B) {
	tracer, stop := startTestTracer(b, WithSampler(NewRateSampler(0)))
	defer stop()

	for n := 0; n < b.N; n++ {
		span := tracer.StartSpan("test")
		span.Finish(StackFrames(64, 0))
	}
}

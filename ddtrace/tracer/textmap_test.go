// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPHeadersCarrierSet(t *testing.T) {
	h := http.Header{}
	c := HTTPHeadersCarrier(h)
	c.Set("A", "x")
	assert.Equal(t, "x", h.Get("A"))
}

func TestHTTPHeadersCarrierForeachKey(t *testing.T) {
	h := http.Header{}
	h.Add("A", "x")
	h.Add("B", "y")
	got := map[string]string{}
	err := HTTPHeadersCarrier(h).ForeachKey(func(k, v string) error {
		got[k] = v
		return nil
	})
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("x", h.Get("A"))
	assert.Equal("y", h.Get("B"))
}

func TestHTTPHeadersCarrierForeachKeyError(t *testing.T) {
	want := errors.New("random error")
	h := http.Header{}
	h.Add("A", "x")
	h.Add("B", "y")
	got := HTTPHeadersCarrier(h).ForeachKey(func(k, v string) error {
		if k == "B" {
			return want
		}
		return nil
	})
	assert.Equal(t, want, got)
}

func TestTextMapCarrierSet(t *testing.T) {
	m := map[string]string{}
	c := TextMapCarrier(m)
	c.Set("a", "b")
	assert.Equal(t, "b", m["a"])
}

func TestTextMapCarrierForeachKey(t *testing.T) {
	want := map[string]string{"A": "x", "B": "y"}
	got := map[string]string{}
	err := TextMapCarrier(want).ForeachKey(func(k, v string) error {
		got[k] = v
		return nil
	})
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(got, want)
}

func TestTextMapCarrierForeachKeyError(t *testing.T) {
	m := map[string]string{"A": "x", "B": "y"}
	want := errors.New("random error")
	got := TextMapCarrier(m).ForeachKey(func(k, v string) error {
		return want
	})
	assert.Equal(t, got, want)
}

func TestW3CExtractsBaggage(t *testing.T) {
	tracer := newUnstartedTracer()
	defer tracer.Stop()
	headers := TextMapCarrier{
		traceparentHeader:      "00-12345678901234567890123456789012-1234567890123456-01",
		tracestateHeader:       "dd=s:2;o:rum;t.usr.id:baz64~~",
		"ot-baggage-something": "someVal",
	}
	s, err := tracer.Extract(headers)
	assert.NoError(t, err)
	found := false
	s.ForeachBaggageItem(func(k, v string) bool {
		if k == "something" {
			found = true
			return false
		}
		return true
	})
	assert.True(t, found)
}

func BenchmarkExtractDatadog(b *testing.B) {
	b.Setenv(headerPropagationStyleExtract, "datadog")
	propagator := NewPropagator(nil)
	carrier := TextMapCarrier(map[string]string{
		DefaultTraceIDHeader:  "1123123132131312313123123",
		DefaultParentIDHeader: "1212321131231312312312312",
		DefaultPriorityHeader: "-1",
		traceTagsHeader: `adad=ada2,adad=ada2,ad1d=ada2,adad=ada2,adad=ada2,
								adad=ada2,adad=aad2,adad=ada2,adad=ada2,adad=ada2,adad=ada2`,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		propagator.Extract(carrier)
	}
}

func BenchmarkExtractW3C(b *testing.B) {
	b.Setenv(headerPropagationStyleExtract, "tracecontext")
	propagator := NewPropagator(nil)
	carrier := TextMapCarrier(map[string]string{
		traceparentHeader: "00-00000000000000001111111111111111-2222222222222222-01",
		tracestateHeader:  "dd=s:2;o:rum;t.dm:-4;t.usr.id:baz64~~,othervendor=t61rcWkgMzE",
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		propagator.Extract(carrier)
	}
}

// Regression test for https://github.com/DataDog/dd-trace-go/issues/1944
func TestPropagatingTagsConcurrency(_ *testing.T) {
	// This test ensures Injection can be done concurrently.
	trc := newUnstartedTracer()
	defer trc.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 1_000; i++ {
		root := trc.StartSpan("test")
		wg.Add(5)
		for i := 0; i < 5; i++ {
			go func() {
				defer wg.Done()
				trc.Inject(root.Context(), TextMapCarrier(make(map[string]string)))
			}()
		}
		wg.Wait()
	}
}

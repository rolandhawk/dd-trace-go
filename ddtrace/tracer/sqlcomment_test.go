// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"testing"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

func FuzzExtract(f *testing.F) {
	testCases := []struct {
		query string
	}{
		{"/*dddbs='whiskey-db',ddps='whiskey-service%20%21%23%24%25%26%27%28%29%2A%2B%2C%2F%3A%3B%3D%3F%40%5B%5D',traceparent='00-0000000000000000<span_id>-<span_id>-00'*/ SELECT * from FOO"},
		{"SELECT * from FOO -- test query"},
		{"/* c */ SELECT traceparent from FOO /**/"},
		{"/*c*/ SELECT traceparent from FOO /**/ /*action='%2Fparam*d',controller='index,'framework='spring',traceparent='<trace-parent>',tracestate='congo%3Dt61rcWkgMzE%2Crojo%3D00f067aa0ba902b7'*/"},
		{"*/ / * * *//*/**/"},
		{""},
	}
	for _, tc := range testCases {
		f.Add(tc.query)
	}
	f.Fuzz(func(t *testing.T, q string) {
		carrier := SQLCommentCarrier{Query: q}
		carrier.Extract() // make sure it doesn't panic
	})
}

func BenchmarkSQLCommentInjection(b *testing.B) {
	tracer, spanCtx, carrier := setupBenchmark()
	defer tracer.Stop()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		carrier.Inject(spanCtx)
	}
}

func BenchmarkSQLCommentExtraction(b *testing.B) {
	tracer, spanCtx, carrier := setupBenchmark()
	defer tracer.Stop()
	carrier.Inject(spanCtx)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		carrier.Extract()
	}
}

func setupBenchmark() (ddtrace.Tracer, ddtrace.SpanContext, SQLCommentCarrier) {
	tracer := newUnstartedTracer(WithService("whiskey-service !#$%&'()*+,/:;=?@[]"), WithEnv("test-env"), WithServiceVersion("1.0.0"))
	root := tracer.StartSpan("service.calling.db", WithSpanID(10))
	root.SetTag(ext.SamplingPriority, 2)
	spanCtx := root.Context()
	carrier := SQLCommentCarrier{Query: "SELECT 1 FROM dual", Mode: DBMPropagationModeFull, DBServiceName: "whiskey-db"}
	return tracer, spanCtx, carrier
}

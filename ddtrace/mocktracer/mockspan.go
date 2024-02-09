// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package mocktracer // import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"

import (
	"fmt"
	"time"

	v2 "github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/internal"
)

// Span is an interface that allows querying a span returned by the mock tracer.
type Span interface {
	// SpanID returns the span's ID.
	SpanID() uint64

	// TraceID returns the span's trace ID.
	TraceID() uint64

	// ParentID returns the span's parent ID.
	ParentID() uint64

	// StartTime returns the time when the span has started.
	StartTime() time.Time

	// FinishTime returns the time when the span has finished.
	FinishTime() time.Time

	// OperationName returns the operation name held by this span.
	OperationName() string

	// Tag returns the value of the tag at key k.
	Tag(k string) interface{}

	// Tags returns a copy of all the tags in this span.
	Tags() map[string]interface{}

	// Context returns the span's SpanContext.
	Context() ddtrace.SpanContext

	// Stringer allows pretty-printing the span's fields for debugging.
	fmt.Stringer
}

type spanV2Adapter struct {
	Span *v2.Span
}

// Context implements Span.
func (sa *spanV2Adapter) Context() ddtrace.SpanContext {
	return &internal.SpanContextV2Adapter{Ctx: sa.Span.Context()}
}

// FinishTime implements Span.
func (sa *spanV2Adapter) FinishTime() time.Time {
	return sa.Span.FinishTime()
}

// OperationName implements Span.
func (sa *spanV2Adapter) OperationName() string {
	return sa.Span.OperationName()
}

// ParentID implements Span.
func (sa *spanV2Adapter) ParentID() uint64 {
	return sa.Span.ParentID()
}

// SpanID implements Span.
func (sa *spanV2Adapter) SpanID() uint64 {
	return sa.Span.SpanID()
}

// StartTime implements Span.
func (sa *spanV2Adapter) StartTime() time.Time {
	return sa.Span.StartTime()
}

// String implements Span.
func (sa *spanV2Adapter) String() string {
	return sa.Span.String()
}

// Tag implements Span.
func (sa *spanV2Adapter) Tag(k string) interface{} {
	return sa.Span.Tag(k)
}

// Tags implements Span.
func (sa *spanV2Adapter) Tags() map[string]interface{} {
	return sa.Span.Tags()
}

// TraceID implements Span.
func (sa *spanV2Adapter) TraceID() uint64 {
	return sa.Span.TraceID()
}

package internal

import (
	"encoding/binary"
	"strconv"

	v2 "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

type SpanContextV1Adapter struct {
	Ctx ddtrace.SpanContext
}

// ForeachBaggageItem implements v2 ddtrace.SpanContext.
func (sca *SpanContextV1Adapter) ForeachBaggageItem(handler func(k string, v string) bool) {
	sca.Ctx.ForeachBaggageItem(handler)
}

// SpanID implements v2 ddtrace.SpanContext.
func (sca *SpanContextV1Adapter) SpanID() uint64 {
	return sca.Ctx.SpanID()
}

// TraceID implements v2 ddtrace.SpanContext.
func (sca *SpanContextV1Adapter) TraceID() string {
	return strconv.FormatUint(sca.Ctx.TraceID(), 10)
}

// TraceIDBytes implements v2 ddtrace.SpanContext.
func (sca *SpanContextV1Adapter) TraceIDBytes() [16]byte {
	var traceIDBytes [16]byte
	binary.BigEndian.PutUint64(traceIDBytes[:], sca.Ctx.TraceID())
	return traceIDBytes
}

// TraceIDLower implements v2 ddtrace.SpanContext.
func (sca *SpanContextV1Adapter) TraceIDLower() uint64 {
	return sca.Ctx.TraceID()
}

type TracerV2Adapter struct {
	Tracer v2.Tracer
}

// Extract implements ddtrace.Tracer.
func (ta *TracerV2Adapter) Extract(carrier interface{}) (ddtrace.SpanContext, error) {
	ctx, err := ta.Tracer.Extract(carrier)
	if err != nil {
		return nil, err
	}
	return &SpanContextV2Adapter{Ctx: ctx}, nil
}

// Inject implements ddtrace.Tracer.
func (ta *TracerV2Adapter) Inject(context ddtrace.SpanContext, carrier interface{}) error {
	ctx := v2.FromGenericCtx(&SpanContextV1Adapter{Ctx: context})
	return ta.Tracer.Inject(ctx, carrier)
}

// StartSpan implements ddtrace.Tracer.
func (ta *TracerV2Adapter) StartSpan(operationName string, opts ...ddtrace.StartSpanOption) ddtrace.Span {
	span := ta.Tracer.StartSpan(operationName, opts...)
	return &SpanV2Adapter{Span: span}
}

// Stop implements ddtrace.Tracer.
func (ta *TracerV2Adapter) Stop() {
	ta.Tracer.Stop()
}

type SpanV2Adapter struct {
	Span *v2.Span
}

// BaggageItem implements ddtrace.Span.
func (sa *SpanV2Adapter) BaggageItem(key string) string {
	return sa.Span.BaggageItem(key)
}

// Context implements ddtrace.Span.
func (sa *SpanV2Adapter) Context() ddtrace.SpanContext {
	Ctx := sa.Span.Context()
	return &SpanContextV2Adapter{Ctx}
}

// Finish implements ddtrace.Span.
func (sa *SpanV2Adapter) Finish(opts ...ddtrace.FinishOption) {
	sa.Span.Finish(opts...)
}

// SetBaggageItem implements ddtrace.Span.
func (sa *SpanV2Adapter) SetBaggageItem(key string, val string) {
	sa.Span.SetBaggageItem(key, val)
}

// SetOperationName implements ddtrace.Span.
func (sa *SpanV2Adapter) SetOperationName(operationName string) {
	sa.Span.SetOperationName(operationName)
}

// SetTag implements ddtrace.Span.
func (sa *SpanV2Adapter) SetTag(key string, value interface{}) {
	sa.Span.SetTag(key, value)
}

func (sa *SpanV2Adapter) Root() ddtrace.Span {
	return &SpanV2Adapter{Span: sa.Span.Root()}
}

type SpanContextV2Adapter struct {
	Ctx *v2.SpanContext
}

// ForeachBaggageItem implements ddtrace.SpanContext.
func (sca *SpanContextV2Adapter) ForeachBaggageItem(handler func(k string, v string) bool) {
	sca.Ctx.ForeachBaggageItem(handler)
}

// SpanID implements ddtrace.SpanContext.
func (sca *SpanContextV2Adapter) SpanID() uint64 {
	return sca.Ctx.SpanID()
}

// TraceID implements ddtrace.SpanContext.
func (sca *SpanContextV2Adapter) TraceID() uint64 {
	return sca.Ctx.TraceIDLower()
}

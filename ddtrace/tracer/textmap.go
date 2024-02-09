// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"net/http"

	v2 "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/internal"
)

// HTTPHeadersCarrier wraps an http.Header as a TextMapWriter and TextMapReader, allowing
// it to be used using the provided Propagator implementation.
type HTTPHeadersCarrier http.Header

var _ TextMapWriter = (*HTTPHeadersCarrier)(nil)
var _ TextMapReader = (*HTTPHeadersCarrier)(nil)

// Set implements TextMapWriter.
func (c HTTPHeadersCarrier) Set(key, val string) {
	http.Header(c).Set(key, val)
}

// ForeachKey implements TextMapReader.
func (c HTTPHeadersCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range c {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// TextMapCarrier allows the use of a regular map[string]string as both TextMapWriter
// and TextMapReader, making it compatible with the provided Propagator.
type TextMapCarrier map[string]string

var _ TextMapWriter = (*TextMapCarrier)(nil)
var _ TextMapReader = (*TextMapCarrier)(nil)

// Set implements TextMapWriter.
func (c TextMapCarrier) Set(key, val string) {
	c[key] = val
}

// ForeachKey conforms to the TextMapReader interface.
func (c TextMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range c {
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

const (
	headerPropagationStyleInject  = "DD_TRACE_PROPAGATION_STYLE_INJECT"
	headerPropagationStyleExtract = "DD_TRACE_PROPAGATION_STYLE_EXTRACT"
	headerPropagationStyle        = "DD_TRACE_PROPAGATION_STYLE"

	headerPropagationStyleInjectDeprecated  = "DD_PROPAGATION_STYLE_INJECT"  // deprecated
	headerPropagationStyleExtractDeprecated = "DD_PROPAGATION_STYLE_EXTRACT" // deprecated
)

const (
	// DefaultBaggageHeaderPrefix specifies the prefix that will be used in
	// HTTP headers or text maps to prefix baggage keys.
	DefaultBaggageHeaderPrefix = "ot-baggage-"

	// DefaultTraceIDHeader specifies the key that will be used in HTTP headers
	// or text maps to store the trace ID.
	DefaultTraceIDHeader = "x-datadog-trace-id"

	// DefaultParentIDHeader specifies the key that will be used in HTTP headers
	// or text maps to store the parent ID.
	DefaultParentIDHeader = "x-datadog-parent-id"

	// DefaultPriorityHeader specifies the key that will be used in HTTP headers
	// or text maps to store the sampling priority value.
	DefaultPriorityHeader = "x-datadog-sampling-priority"
)

// originHeader specifies the name of the header indicating the origin of the trace.
// It is used with the Synthetics product and usually has the value "synthetics".
const originHeader = "x-datadog-origin"

// traceTagsHeader holds the propagated trace tags
const traceTagsHeader = "x-datadog-tags"

// propagationExtractMaxSize limits the total size of incoming propagated tags to parse
const propagationExtractMaxSize = 512

// PropagatorConfig defines the configuration for initializing a propagator.
type PropagatorConfig = v2.PropagatorConfig

type propagatorV1Adapter struct {
	Propagator Propagator
}

// Extract implements tracer.Propagator.
func (pa *propagatorV1Adapter) Extract(carrier interface{}) (*v2.SpanContext, error) {
	ctx, err := pa.Propagator.Extract(carrier)
	if err != nil {
		return nil, err
	}
	return v2.FromGenericCtx(&internal.SpanContextV1Adapter{Ctx: ctx}), nil
}

// Inject implements tracer.Propagator.
func (pa *propagatorV1Adapter) Inject(context *v2.SpanContext, carrier interface{}) error {
	ctx := &internal.SpanContextV2Adapter{Ctx: context}
	return pa.Propagator.Inject(ctx, carrier)
}

type propagatorV2Adapter struct {
	Propagator v2.Propagator
}

// Extract implements Propagator.
func (pa *propagatorV2Adapter) Extract(carrier interface{}) (ddtrace.SpanContext, error) {
	ctx, err := pa.Propagator.Extract(carrier)
	if err != nil {
		return nil, err
	}
	return &internal.SpanContextV2Adapter{Ctx: ctx}, nil
}

// Inject implements Propagator.
func (pa *propagatorV2Adapter) Inject(context ddtrace.SpanContext, carrier interface{}) error {
	ctx := v2.FromGenericCtx(&internal.SpanContextV1Adapter{Ctx: context})
	return pa.Propagator.Inject(ctx, carrier)
}

// NewPropagator returns a new propagator which uses TextMap to inject
// and extract values. It propagates trace and span IDs and baggage.
// To use the defaults, nil may be provided in place of the config.
//
// The inject and extract propagators are determined using environment variables
// with the following order of precedence:
//  1. DD_TRACE_PROPAGATION_STYLE_INJECT
//  2. DD_PROPAGATION_STYLE_INJECT (deprecated)
//  3. DD_TRACE_PROPAGATION_STYLE (applies to both inject and extract)
//  4. If none of the above, use default values
func NewPropagator(cfg *PropagatorConfig, propagators ...Propagator) Propagator {
	converted := make([]v2.Propagator, len(propagators))
	for i, p := range propagators {
		converted[i] = &propagatorV1Adapter{
			Propagator: p,
		}
	}
	p := v2.NewPropagator(cfg, converted...)
	return &propagatorV2Adapter{
		Propagator: p,
	}
}

const (
	traceparentHeader = "traceparent"
	tracestateHeader  = "tracestate"
)

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	v2 "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/internal"
)

// Sampler is the generic interface of any sampler. It must be safe for concurrent use.
type Sampler interface {
	// Sample returns true if the given span should be sampled.
	Sample(span Span) bool
}

// RateSampler is a sampler implementation which randomly selects spans using a
// provided rate. For example, a rate of 0.75 will permit 75% of the spans.
// RateSampler implementations should be safe for concurrent use.
type RateSampler interface {
	Sampler

	// Rate returns the current sample rate.
	Rate() float64

	// SetRate sets a new sample rate.
	SetRate(rate float64)
}

type rateSamplerV2Adapter struct {
	Sampler v2.RateSampler
}

// Rate implements RateSampler.
func (rsa *rateSamplerV2Adapter) Rate() float64 {
	return rsa.Sampler.Rate()
}

// Sample implements RateSampler.
func (rsa *rateSamplerV2Adapter) Sample(span ddtrace.Span) bool {
	sp := span.(*internal.SpanV2Adapter)
	return rsa.Sampler.Sample(sp.Span)
}

// SetRate implements RateSampler.
func (rsa *rateSamplerV2Adapter) SetRate(rate float64) {
	rsa.Sampler.SetRate(rate)
}

// NewAllSampler is a short-hand for NewRateSampler(1). It is all-permissive.
func NewAllSampler() RateSampler { return NewRateSampler(1) }

// NewRateSampler returns an initialized RateSampler with a given sample rate.
func NewRateSampler(rate float64) RateSampler {
	return &rateSamplerV2Adapter{
		Sampler: v2.NewRateSampler(rate),
	}
}

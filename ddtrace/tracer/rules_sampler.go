// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"regexp"

	v2 "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

// SamplingRule is used for applying sampling rates to spans that match
// the service name, operation name or both.
// For basic usage, consider using the helper functions ServiceRule, NameRule, etc.
type SamplingRule = v2.SamplingRule

// SamplingRuleType represents a type of sampling rule spans are matched against.
type SamplingRuleType = v2.SamplingRuleType

// ServiceRule returns a SamplingRule that applies the provided sampling rate
// to spans that match the service name provided.
func ServiceRule(service string, rate float64) SamplingRule {
	r := v2.TraceSamplingRules(v2.Rule{
		ServiceGlob: service,
		Rate:        rate,
	})
	return r[0]
}

// NameRule returns a SamplingRule that applies the provided sampling rate
// to spans that match the operation name provided.
func NameRule(name string, rate float64) SamplingRule {
	r := v2.TraceSamplingRules(v2.Rule{
		NameGlob: name,
		Rate:     rate,
	})
	return r[0]
}

// NameServiceRule returns a SamplingRule that applies the provided sampling rate
// to spans matching both the operation and service names provided.
func NameServiceRule(name string, service string, rate float64) SamplingRule {
	r := v2.TraceSamplingRules(v2.Rule{
		ServiceGlob: service,
		NameGlob:    name,
		Rate:        rate,
	})
	return r[0]
}

// RateRule returns a SamplingRule that applies the provided sampling rate to all spans.
func RateRule(rate float64) SamplingRule {
	return SamplingRule{
		Rate: rate,
	}
}

// TagsResourceRule returns a SamplingRule that applies the provided sampling rate to traces with spans that match
// resource, name, service and tags provided.
func TagsResourceRule(tags map[string]*regexp.Regexp, resource, name, service string, rate float64) SamplingRule {
	converted := make(map[string]string, len(tags))
	for k, v := range tags {
		converted[k] = v.String()
	}
	r := v2.TraceSamplingRules(v2.Rule{
		ServiceGlob:  service,
		NameGlob:     name,
		ResourceGlob: resource,
		Rate:         rate,
		Tags:         converted,
	})
	return r[0]
}

// SpanTagsResourceRule returns a SamplingRule that applies the provided sampling rate to spans that match
// resource, name, service and tags provided. Values of the tags map are expected to be in glob format.
func SpanTagsResourceRule(tags map[string]string, resource, name, service string, rate float64) SamplingRule {
	r := v2.SpanSamplingRules(v2.Rule{
		ServiceGlob:  service,
		NameGlob:     name,
		ResourceGlob: resource,
		Rate:         rate,
		Tags:         tags,
	})
	return r[0]
}

// SpanNameServiceRule returns a SamplingRule of type SamplingRuleSpan that applies
// the provided sampling rate to all spans matching the operation and service name glob patterns provided.
// Operation and service fields must be valid glob patterns.
func SpanNameServiceRule(name, service string, rate float64) SamplingRule {
	r := v2.SpanSamplingRules(v2.Rule{
		ServiceGlob:  service,
		NameGlob:     name,
		Rate:         rate,
		MaxPerSecond: 0,
	})
	return r[0]
}

// SpanNameServiceMPSRule returns a SamplingRule of type SamplingRuleSpan that applies
// the provided sampling rate to all spans matching the operation and service name glob patterns
// up to the max number of spans per second that can be sampled.
// Operation and service fields must be valid glob patterns.
func SpanNameServiceMPSRule(name, service string, rate, limit float64) SamplingRule {
	r := v2.SpanSamplingRules(v2.Rule{
		ServiceGlob:  service,
		NameGlob:     name,
		Rate:         rate,
		MaxPerSecond: limit,
	})
	return r[0]
}

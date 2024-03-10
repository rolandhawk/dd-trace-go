// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package gocql

import (
	v2 "github.com/DataDog/dd-trace-go/v2/contrib/gocql/gocql"
)

// WrapOption represents an option that can be passed to WrapQuery.
type WrapOption = v2.WrapOption

// WithServiceName sets the given service name for the returned query.
func WithServiceName(name string) WrapOption {
	return v2.WithService(name)
}

// WithResourceName sets a custom resource name to be used with the traced query.
// By default, the query statement is extracted automatically. This method should
// be used when a different resource name is desired or in performance critical
// environments. The gocql library returns the query statement using an fmt.Sprintf
// call, which can be costly when called repeatedly. Using WithResourceName will
// avoid that call. Under normal circumstances, it is safe to rely on the default.
func WithResourceName(name string) WrapOption {
	return v2.WithResourceName(name)
}

// WithAnalytics enables Trace Analytics for all started spans.
func WithAnalytics(on bool) WrapOption {
	return v2.WithAnalytics(on)
}

// WithAnalyticsRate sets the sampling rate for Trace Analytics events
// correlated to started spans.
func WithAnalyticsRate(rate float64) WrapOption {
	return v2.WithAnalyticsRate(rate)
}

// NoDebugStack prevents stack traces from being attached to spans finishing
// with an error. This is useful in situations where errors are frequent and
// performance is critical.
func NoDebugStack() WrapOption {
	return v2.NoDebugStack()
}

// WithErrorCheck specifies a function fn which determines whether the passed
// error should be marked as an error. The fn is called whenever a CQL request
// finishes with an error.
func WithErrorCheck(fn func(err error) bool) WrapOption {
	return v2.WithErrorCheck(fn)
}

// WithCustomTag will attach the value to the span tagged by the key.
func WithCustomTag(key string, value interface{}) WrapOption {
	return v2.WithCustomTag(key, value)
}

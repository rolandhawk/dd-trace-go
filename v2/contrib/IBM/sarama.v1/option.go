// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package sarama

import (
	"math"

	"github.com/DataDog/dd-trace-go/v2/internal"
	"github.com/DataDog/dd-trace-go/v2/internal/namingschema"
)

const defaultServiceName = "kafka"

type config struct {
	consumerServiceName string
	producerServiceName string
	consumerSpanName    string
	producerSpanName    string
	analyticsRate       float64
}

func defaults(cfg *config) {
	cfg.consumerServiceName = namingschema.NewDefaultServiceName(defaultServiceName).GetName()
	cfg.producerServiceName = namingschema.NewDefaultServiceName(
		defaultServiceName,
		namingschema.WithOverrideV0(defaultServiceName),
	).GetName()

	cfg.consumerSpanName = namingschema.NewKafkaInboundOp().GetName()
	cfg.producerSpanName = namingschema.NewKafkaOutboundOp().GetName()

	// cfg.analyticsRate = globalconfig.AnalyticsRate()
	if internal.BoolEnv("DD_TRACE_SARAMA_ANALYTICS_ENABLED", false) {
		cfg.analyticsRate = 1.0
	} else {
		cfg.analyticsRate = math.NaN()
	}
}

// Option describes options for the Sarama integration.
type Option interface {
	apply(*config)
}

// OptionFn represents options applicable to WrapConsumer, WrapPartitionConsumer, WrapAsyncProducer and WrapSyncProducer.
type OptionFn func(*config)

func (fn OptionFn) apply(cfg *config) {
	fn(cfg)
}

// WithService sets the given service name for the intercepted client.
func WithService(name string) OptionFn {
	return func(cfg *config) {
		cfg.consumerServiceName = name
		cfg.producerServiceName = name
	}
}

// WithAnalytics enables Trace Analytics for all started spans.
func WithAnalytics(on bool) OptionFn {
	return func(cfg *config) {
		if on {
			cfg.analyticsRate = 1.0
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

// WithAnalyticsRate sets the sampling rate for Trace Analytics events
// correlated to started spans.
func WithAnalyticsRate(rate float64) OptionFn {
	return func(cfg *config) {
		if rate >= 0.0 && rate <= 1.0 {
			cfg.analyticsRate = rate
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

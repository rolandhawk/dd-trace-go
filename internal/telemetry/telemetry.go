// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.

// Package telemetry implements a client for sending telemetry information to
// Datadog regarding usage of an APM library such as tracing or profiling.
package telemetry

import (
	"github.com/DataDog/dd-trace-go/v2/v1internal/telemetry"
)

// Integrations returns which integrations are tracked by telemetry.
func Integrations() []Integration {
	return telemetry.Integrations()
}

// LoadIntegration notifies telemetry that an integration is being used.
func LoadIntegration(name string) {
	telemetry.LoadIntegration(name)
}

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package mocktracer

import (
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

func TestStart(t *testing.T) {
	trc := Start()
	if _, ok := trc.(tracer.Tracer); !ok {
		t.Fail()
	}
	if tt, ok := tracer.GetGlobalTracer().(Tracer); !ok || tt != trc {
		t.Fail()
	}
}

func TestTracerStop(t *testing.T) {
	Start().Stop()
	if _, ok := tracer.GetGlobalTracer().(*tracer.NoopTracer); !ok {
		t.Fail()
	}
}

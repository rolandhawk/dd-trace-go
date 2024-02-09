// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/internal"
)

// SQLCommentInjectionMode represents the mode of SQL comment injection.
//
// Deprecated: Use DBMPropagationMode instead.
type SQLCommentInjectionMode DBMPropagationMode

const (
	// SQLInjectionUndefined represents the comment injection mode is not set. This is the same as SQLInjectionDisabled.
	SQLInjectionUndefined SQLCommentInjectionMode = SQLCommentInjectionMode(DBMPropagationModeUndefined)
	// SQLInjectionDisabled represents the comment injection mode where all injection is disabled.
	SQLInjectionDisabled SQLCommentInjectionMode = SQLCommentInjectionMode(DBMPropagationModeDisabled)
	// SQLInjectionModeService represents the comment injection mode where only service tags (name, env, version) are injected.
	SQLInjectionModeService SQLCommentInjectionMode = SQLCommentInjectionMode(DBMPropagationModeService)
	// SQLInjectionModeFull represents the comment injection mode where both service tags and tracing tags. Tracing tags include span id, trace id and sampling priority.
	SQLInjectionModeFull SQLCommentInjectionMode = SQLCommentInjectionMode(DBMPropagationModeFull)
)

// DBMPropagationMode represents the mode of dbm propagation.
//
// Note that enabling sql comment propagation results in potentially confidential data (service names)
// being stored in the databases which can then be accessed by other 3rd parties that have been granted
// access to the database.
type DBMPropagationMode string

const (
	// DBMPropagationModeUndefined represents the dbm propagation mode not being set. This is the same as DBMPropagationModeDisabled.
	DBMPropagationModeUndefined DBMPropagationMode = ""
	// DBMPropagationModeDisabled represents the dbm propagation mode where all propagation is disabled.
	DBMPropagationModeDisabled DBMPropagationMode = "disabled"
	// DBMPropagationModeService represents the dbm propagation mode where only service tags (name, env, version) are propagated to dbm.
	DBMPropagationModeService DBMPropagationMode = "service"
	// DBMPropagationModeFull represents the dbm propagation mode where both service tags and tracing tags are propagated. Tracing tags include span id, trace id and the sampled flag.
	DBMPropagationModeFull DBMPropagationMode = "full"
)

// Key names for SQL comment tags.
const (
	sqlCommentTraceParent   = "traceparent"
	sqlCommentParentService = "ddps"
	sqlCommentDBService     = "dddbs"
	sqlCommentParentVersion = "ddpv"
	sqlCommentEnv           = "dde"
)

// Current trace context version (see https://www.w3.org/TR/trace-context/#version)
const w3cContextVersion = "00"

// SQLCommentCarrier is a carrier implementation that injects a span context in a SQL query in the form
// of a sqlcommenter formatted comment prepended to the original query text.
// See https://google.github.io/sqlcommenter/spec/ for more details.
type SQLCommentCarrier struct {
	Query         string
	Mode          DBMPropagationMode
	DBServiceName string
	SpanID        uint64
	carrier       *tracer.SQLCommentCarrier
}

// Inject injects a span context in the carrier's Query field as a comment.
func (c *SQLCommentCarrier) Inject(spanCtx ddtrace.SpanContext) error {
	if c.carrier == nil {
		c.carrier = &tracer.SQLCommentCarrier{
			Query:         c.Query,
			Mode:          tracer.DBMPropagationMode(c.Mode),
			DBServiceName: c.DBServiceName,
			SpanID:        c.SpanID,
		}
	}
	return c.carrier.Inject(
		tracer.FromGenericCtx(&internal.SpanContextV1Adapter{
			Ctx: spanCtx,
		}),
	)
}

// Extract parses for key value attributes in a sql query injected with trace information in order to build a span context
func (c *SQLCommentCarrier) Extract() (ddtrace.SpanContext, error) {
	if c.carrier == nil {
		c.carrier = &tracer.SQLCommentCarrier{
			Query:         c.Query,
			Mode:          tracer.DBMPropagationMode(c.Mode),
			DBServiceName: c.DBServiceName,
			SpanID:        c.SpanID,
		}
	}
	ctx, err := c.carrier.Extract()
	if err != nil {
		return nil, err
	}
	return &internal.SpanContextV2Adapter{Ctx: ctx}, nil
}

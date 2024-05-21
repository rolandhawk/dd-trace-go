// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

// Package chi provides tracing functions for tracing the go-chi/chi/v5 package (https://github.com/go-chi/chi).
package chi // import "github.com/DataDog/dd-trace-go/contrib/go-chi/chi.v5"/v2

import (
	"fmt"
	"math"
	"net/http"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/DataDog/dd-trace-go/v2/internal/appsec"
	"github.com/DataDog/dd-trace-go/v2/internal/contrib/httptrace"
	"github.com/DataDog/dd-trace-go/v2/internal/contrib/options"
	"github.com/DataDog/dd-trace-go/v2/internal/log"
	"github.com/DataDog/dd-trace-go/v2/internal/telemetry"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const componentName = "go-chi/chi.v5"

func init() {
	telemetry.LoadIntegration(componentName)
	tracer.MarkIntegrationImported("github.com/go-chi/chi/v5")
}

// Middleware returns middleware that will trace incoming requests.
func Middleware(opts ...Option) func(next http.Handler) http.Handler {
	cfg := new(config)
	defaults(cfg)
	for _, fn := range opts {
		fn.apply(cfg)
	}
	log.Debug("contrib/go-chi/chi.v5: Configuring Middleware: %#v", cfg)
	spanOpts := append(cfg.spanOpts, tracer.ServiceName(cfg.serviceName),
		tracer.Tag(ext.Component, componentName),
		tracer.Tag(ext.SpanKind, ext.SpanKindServer))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.ignoreRequest(r) {
				next.ServeHTTP(w, r)
				return
			}
			opts := options.Expand(spanOpts, 0, 2) // opts must be a copy of spanOpts, locally scoped, to avoid races.
			if !math.IsNaN(cfg.analyticsRate) {
				opts = append(opts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
			}
			opts = append(opts, httptrace.HeaderTagsFromRequest(r, cfg.headerTags))
			span, ctx := httptrace.StartRequestSpan(r, opts...)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				status := ww.Status()
				var opts []tracer.FinishOption
				if cfg.isStatusError(status) {
					opts = []tracer.FinishOption{tracer.WithError(fmt.Errorf("%d: %s", status, http.StatusText(status)))}
				}
				httptrace.FinishRequestSpan(span, status, opts...)
			}()

			// pass the span through the request context
			r = r.WithContext(ctx)

			next := next // avoid modifying the value of next in the outer closure scope
			if appsec.Enabled() && !cfg.appsecDisabled {
				next = withAppsec(next, r, span, &cfg.appsecConfig)
				// Note that the following response writer passed to the handler
				// implements the `interface { Status() int }` expected by httpsec.
			}

			// pass the span through the request context and serve the request to the next middleware
			next.ServeHTTP(ww, r)

			routePattern := cfg.modifyResourceName(chi.RouteContext(r.Context()).RoutePattern())
			span.SetTag(ext.HTTPRoute, routePattern)
			var resourceName string
			if cfg.resourceNamer != nil {
				resourceName = cfg.resourceNamer(r)
			} else {
				resourceName = routePattern
				if resourceName == "" {
					resourceName = "unknown"
				}
				resourceName = r.Method + " " + resourceName
			}
			span.SetTag(ext.ResourceName, resourceName)
		})
	}
}

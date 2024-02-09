// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

package tracer

import (
	"context"
	"net/http"
	"strings"
	"time"

	v2 "github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/internal"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/globalconfig"
	"gopkg.in/DataDog/dd-trace-go.v1/internal/normalizer"
)

var contribIntegrations = map[string]struct {
	name     string // user readable name for startup logs
	imported bool   // true if the user has imported the integration
}{
	"github.com/99designs/gqlgen":                   {"gqlgen", false},
	"github.com/aws/aws-sdk-go":                     {"AWS SDK", false},
	"github.com/aws/aws-sdk-go-v2":                  {"AWS SDK v2", false},
	"github.com/bradfitz/gomemcache":                {"Memcache", false},
	"cloud.google.com/go/pubsub.v1":                 {"Pub/Sub", false},
	"github.com/confluentinc/confluent-kafka-go":    {"Kafka (confluent)", false},
	"github.com/confluentinc/confluent-kafka-go/v2": {"Kafka (confluent) v2", false},
	"database/sql":                                  {"SQL", false},
	"github.com/dimfeld/httptreemux/v5":             {"HTTP Treemux", false},
	"github.com/elastic/go-elasticsearch/v6":        {"Elasticsearch v6", false},
	"github.com/emicklei/go-restful":                {"go-restful", false},
	"github.com/emicklei/go-restful/v3":             {"go-restful v3", false},
	"github.com/garyburd/redigo":                    {"Redigo (dep)", false},
	"github.com/gin-gonic/gin":                      {"Gin", false},
	"github.com/globalsign/mgo":                     {"MongoDB (mgo)", false},
	"github.com/go-chi/chi":                         {"chi", false},
	"github.com/go-chi/chi/v5":                      {"chi v5", false},
	"github.com/go-pg/pg/v10":                       {"go-pg v10", false},
	"github.com/go-redis/redis":                     {"Redis", false},
	"github.com/go-redis/redis/v7":                  {"Redis v7", false},
	"github.com/go-redis/redis/v8":                  {"Redis v8", false},
	"go.mongodb.org/mongo-driver":                   {"MongoDB", false},
	"github.com/gocql/gocql":                        {"Cassandra", false},
	"github.com/gofiber/fiber/v2":                   {"Fiber", false},
	"github.com/gomodule/redigo":                    {"Redigo", false},
	"google.golang.org/api":                         {"Google API", false},
	"google.golang.org/grpc":                        {"gRPC", false},
	"google.golang.org/grpc/v12":                    {"gRPC v12", false},
	"gopkg.in/jinzhu/gorm.v1":                       {"Gorm (gopkg)", false},
	"github.com/gorilla/mux":                        {"Gorilla Mux", false},
	"gorm.io/gorm.v1":                               {"Gorm v1", false},
	"github.com/graph-gophers/graphql-go":           {"GraphQL", false},
	"github.com/hashicorp/consul/api":               {"Consul", false},
	"github.com/hashicorp/vault/api":                {"Vault", false},
	"github.com/jinzhu/gorm":                        {"Gorm", false},
	"github.com/jmoiron/sqlx":                       {"SQLx", false},
	"github.com/julienschmidt/httprouter":           {"HTTP Router", false},
	"k8s.io/client-go/kubernetes":                   {"Kubernetes", false},
	"github.com/labstack/echo":                      {"echo", false},
	"github.com/labstack/echo/v4":                   {"echo v4", false},
	"github.com/miekg/dns":                          {"miekg/dns", false},
	"net/http":                                      {"HTTP", false},
	"gopkg.in/olivere/elastic.v5":                   {"Elasticsearch v5", false},
	"gopkg.in/olivere/elastic.v3":                   {"Elasticsearch v3", false},
	"github.com/redis/go-redis/v9":                  {"Redis v9", false},
	"github.com/segmentio/kafka-go":                 {"Kafka v0", false},
	"github.com/IBM/sarama":                         {"IBM sarama", false},
	"github.com/Shopify/sarama":                     {"Shopify sarama", false},
	"github.com/sirupsen/logrus":                    {"Logrus", false},
	"github.com/syndtr/goleveldb":                   {"LevelDB", false},
	"github.com/tidwall/buntdb":                     {"BuntDB", false},
	"github.com/twitchtv/twirp":                     {"Twirp", false},
	"github.com/urfave/negroni":                     {"Negroni", false},
	"github.com/valyala/fasthttp":                   {"FastHTTP", false},
	"github.com/zenazn/goji":                        {"Goji", false},
}

// StartOption represents a function that can be provided as a parameter to Start.
type StartOption = v2.StartOption

type integrationConfig struct {
	Instrumented bool   `json:"instrumented"`      // indicates if the user has imported and used the integration
	Available    bool   `json:"available"`         // indicates if the user is using a library that can be used with DataDog integrations
	Version      string `json:"available_version"` // if available, indicates the version of the library the user has
}

// agentFeatures holds information about the trace-agent's capabilities.
// When running WithLambdaMode, a zero-value of this struct will be used
// as features.
type agentFeatures struct {
	// DropP0s reports whether it's ok for the tracer to not send any
	// P0 traces to the agent.
	DropP0s bool

	// Stats reports whether the agent can receive client-computed stats on
	// the /v0.6/stats endpoint.
	Stats bool

	// DataStreams reports whether the agent can receive data streams stats on
	// the /v0.1/pipeline_stats endpoint.
	DataStreams bool

	// StatsdPort specifies the Dogstatsd port as provided by the agent.
	// If it's the default, it will be 0, which means 8125.
	StatsdPort int

	// featureFlags specifies all the feature flags reported by the trace-agent.
	featureFlags map[string]struct{}
}

// HasFlag reports whether the agent has set the feat feature flag.
func (a *agentFeatures) HasFlag(feat string) bool {
	_, ok := a.featureFlags[feat]
	return ok
}

// MarkIntegrationImported labels the given integration as imported
func MarkIntegrationImported(integration string) bool {
	return v2.MarkIntegrationImported(integration)
}

// WithFeatureFlags specifies a set of feature flags to enable. Please take into account
// that most, if not all features flags are considered to be experimental and result in
// unexpected bugs.
func WithFeatureFlags(feats ...string) StartOption {
	return v2.WithFeatureFlags(feats...)
}

// WithLogger sets logger as the tracer's error printer.
// Diagnostic and startup tracer logs are prefixed to simplify the search within logs.
// If JSON logging format is required, it's possible to wrap tracer logs using an existing JSON logger with this
// function. To learn more about this possibility, please visit: https://github.com/DataDog/dd-trace-go/issues/2152#issuecomment-1790586933
func WithLogger(logger ddtrace.Logger) StartOption {
	return v2.WithLogger(logger)
}

// WithPrioritySampling is deprecated, and priority sampling is enabled by default.
// When using distributed tracing, the priority sampling value is propagated in order to
// get all the parts of a distributed trace sampled.
// To learn more about priority sampling, please visit:
// https://docs.datadoghq.com/tracing/getting_further/trace_sampling_and_storage/#priority-sampling-for-distributed-tracing
func WithPrioritySampling() StartOption {
	return nil
}

// WithDebugStack can be used to globally enable or disable the collection of stack traces when
// spans finish with errors. It is enabled by default. This is a global version of the NoDebugStack
// FinishOption.
func WithDebugStack(enabled bool) StartOption {
	return WithDebugStack(enabled)
}

// WithDebugMode enables debug mode on the tracer, resulting in more verbose logging.
func WithDebugMode(enabled bool) StartOption {
	return v2.WithDebugMode(enabled)
}

// WithLambdaMode enables lambda mode on the tracer, for use with AWS Lambda.
// This option is only required if the the Datadog Lambda Extension is not
// running.
func WithLambdaMode(enabled bool) StartOption {
	return v2.WithLambdaMode(enabled)
}

// WithSendRetries enables re-sending payloads that are not successfully
// submitted to the agent.  This will cause the tracer to retry the send at
// most `retries` times.
func WithSendRetries(retries int) StartOption {
	return v2.WithSendRetries(retries)
}

// WithPropagator sets an alternative propagator to be used by the tracer.
func WithPropagator(p Propagator) StartOption {
	return v2.WithPropagator(&propagatorV1Adapter{Propagator: p})
}

// WithServiceName is deprecated. Please use WithService.
// If you are using an older version and you are upgrading from WithServiceName
// to WithService, please note that WithService will determine the service name of
// server and framework integrations.
func WithServiceName(name string) StartOption {
	return v2.WithService(name)
}

// WithService sets the default service name for the program.
func WithService(name string) StartOption {
	return v2.WithService(name)
}

// WithGlobalServiceName causes contrib libraries to use the global service name and not any locally defined service name.
// This is synonymous with `DD_TRACE_REMOVE_INTEGRATION_SERVICE_NAMES_ENABLED`.
func WithGlobalServiceName(enabled bool) StartOption {
	return v2.WithGlobalServiceName(enabled)
}

// WithAgentAddr sets the address where the agent is located. The default is
// localhost:8126. It should contain both host and port.
func WithAgentAddr(addr string) StartOption {
	return v2.WithAgentAddr(addr)
}

// WithEnv sets the environment to which all traces started by the tracer will be submitted.
// The default value is the environment variable DD_ENV, if it is set.
func WithEnv(env string) StartOption {
	return v2.WithEnv(env)
}

// WithServiceMapping determines service "from" to be renamed to service "to".
// This option is is case sensitive and can be used multiple times.
func WithServiceMapping(from, to string) StartOption {
	return v2.WithServiceMapping(from, to)
}

// WithPeerServiceDefaults sets default calculation for peer.service.
func WithPeerServiceDefaults(enabled bool) StartOption {
	// TODO: add link to public docs
	return v2.WithPeerServiceDefaults(enabled)
}

// WithPeerServiceMapping determines the value of the peer.service tag "from" to be renamed to service "to".
func WithPeerServiceMapping(from, to string) StartOption {
	return v2.WithPeerServiceMapping(from, to)
}

// WithGlobalTag sets a key/value pair which will be set as a tag on all spans
// created by tracer. This option may be used multiple times.
func WithGlobalTag(k string, v interface{}) StartOption {
	return v2.WithGlobalTag(k, v)
}

// WithSampler sets the given sampler to be used with the tracer. By default
// an all-permissive sampler is used.
func WithSampler(s Sampler) StartOption {
	rules := []SamplingRule{}
	// TODO (darccio): is it possible to map a Sampler to a set of rules?
	return v2.WithSamplingRules(rules)
}

// WithHTTPRoundTripper is deprecated. Please consider using WithHTTPClient instead.
// The function allows customizing the underlying HTTP transport for emitting spans.
func WithHTTPRoundTripper(r http.RoundTripper) StartOption {
	return v2.WithHTTPClient(&http.Client{Transport: r})
}

// WithHTTPClient specifies the HTTP client to use when emitting spans to the agent.
func WithHTTPClient(client *http.Client) StartOption {
	return v2.WithHTTPClient(client)
}

// WithUDS configures the HTTP client to dial the Datadog Agent via the specified Unix Domain Socket path.
func WithUDS(socketPath string) StartOption {
	return v2.WithUDS(socketPath)
}

// WithAnalytics allows specifying whether Trace Search & Analytics should be enabled
// for integrations.
func WithAnalytics(on bool) StartOption {
	return v2.WithAnalytics(on)
}

// WithAnalyticsRate sets the global sampling rate for sampling APM events.
func WithAnalyticsRate(rate float64) StartOption {
	return v2.WithAnalyticsRate(rate)
}

// WithRuntimeMetrics enables automatic collection of runtime metrics every 10 seconds.
func WithRuntimeMetrics() StartOption {
	return v2.WithRuntimeMetrics()
}

// WithDogstatsdAddress specifies the address to connect to for sending metrics to the Datadog
// Agent. It should be a "host:port" string, or the path to a unix domain socket.If not set, it
// attempts to determine the address of the statsd service according to the following rules:
//  1. Look for /var/run/datadog/dsd.socket and use it if present. IF NOT, continue to #2.
//  2. The host is determined by DD_AGENT_HOST, and defaults to "localhost"
//  3. The port is retrieved from the agent. If not present, it is determined by DD_DOGSTATSD_PORT, and defaults to 8125
//
// This option is in effect when WithRuntimeMetrics is enabled.
func WithDogstatsdAddress(addr string) StartOption {
	return v2.WithDogstatsdAddress(addr)
}

// WithSamplingRules specifies the sampling rates to apply to spans based on the
// provided rules.
func WithSamplingRules(rules []SamplingRule) StartOption {
	return v2.WithSamplingRules(rules)
}

// WithServiceVersion specifies the version of the service that is running. This will
// be included in spans from this service in the "version" tag, provided that
// span service name and config service name match. Do NOT use with WithUniversalVersion.
func WithServiceVersion(version string) StartOption {
	return v2.WithServiceVersion(version)
}

// WithUniversalVersion specifies the version of the service that is running, and will be applied to all spans,
// regardless of whether span service name and config service name match.
// See: WithService, WithServiceVersion. Do NOT use with WithServiceVersion.
func WithUniversalVersion(version string) StartOption {
	return v2.WithUniversalVersion(version)
}

// WithHostname allows specifying the hostname with which to mark outgoing traces.
func WithHostname(name string) StartOption {
	return v2.WithHostname(name)
}

// WithTraceEnabled allows specifying whether tracing will be enabled
func WithTraceEnabled(enabled bool) StartOption {
	return v2.WithTraceEnabled(enabled)
}

// WithLogStartup allows enabling or disabling the startup log.
func WithLogStartup(enabled bool) StartOption {
	return v2.WithLogStartup(enabled)
}

// WithProfilerCodeHotspots enables the code hotspots integration between the
// tracer and profiler. This is done by automatically attaching pprof labels
// called "span id" and "local root span id" when new spans are created. You
// should not use these label names in your own code when this is enabled. The
// enabled value defaults to the value of the
// DD_PROFILING_CODE_HOTSPOTS_COLLECTION_ENABLED env variable or true.
func WithProfilerCodeHotspots(enabled bool) StartOption {
	return v2.WithProfilerCodeHotspots(enabled)
}

// WithProfilerEndpoints enables the endpoints integration between the tracer
// and profiler. This is done by automatically attaching a pprof label called
// "trace endpoint" holding the resource name of the top-level service span if
// its type is "http", "rpc" or "" (default). You should not use this label
// name in your own code when this is enabled. The enabled value defaults to
// the value of the DD_PROFILING_ENDPOINT_COLLECTION_ENABLED env variable or
// true.
func WithProfilerEndpoints(enabled bool) StartOption {
	return v2.WithProfilerEndpoints(enabled)
}

// WithDebugSpansMode enables debugging old spans that may have been
// abandoned, which may prevent traces from being set to the Datadog
// Agent, especially if partial flushing is off.
// This setting can also be configured by setting DD_TRACE_DEBUG_ABANDONED_SPANS
// to true. The timeout will default to 10 minutes, unless overwritten
// by DD_TRACE_ABANDONED_SPAN_TIMEOUT.
// This feature is disabled by default. Turning on this debug mode may
// be expensive, so it should only be enabled for debugging purposes.
func WithDebugSpansMode(timeout time.Duration) StartOption {
	return v2.WithDebugSpansMode(timeout)
}

// WithPartialFlushing enables flushing of partially finished traces.
// This is done after "numSpans" have finished in a single local trace at
// which point all finished spans in that trace will be flushed, freeing up
// any memory they were consuming. This can also be configured by setting
// DD_TRACE_PARTIAL_FLUSH_ENABLED to true, which will default to 1000 spans
// unless overriden with DD_TRACE_PARTIAL_FLUSH_MIN_SPANS. Partial flushing
// is disabled by default.
func WithPartialFlushing(numSpans int) StartOption {
	return v2.WithPartialFlushing(numSpans)
}

// WithStatsComputation enables client-side stats computation, allowing
// the tracer to compute stats from traces. This can reduce network traffic
// to the Datadog Agent, and produce more accurate stats data.
// This can also be configured by setting DD_TRACE_STATS_COMPUTATION_ENABLED to true.
// Client-side stats is off by default.
func WithStatsComputation(enabled bool) StartOption {
	return v2.WithStatsComputation(enabled)
}

// WithOrchestrion configures Orchestrion's auto-instrumentation metadata.
// This option is only intended to be used by Orchestrion https://github.com/DataDog/orchestrion
func WithOrchestrion(metadata map[string]string) StartOption {
	return v2.WithOrchestrion(metadata)
}

// StartSpanOption is a configuration option for StartSpan. It is aliased in order
// to help godoc group all the functions returning it together. It is considered
// more correct to refer to it as the type as the origin, ddtrace.StartSpanOption.
type StartSpanOption = ddtrace.StartSpanOption

// Tag sets the given key/value pair as a tag on the started Span.
func Tag(k string, v interface{}) StartSpanOption {
	return func(cfg *v2.StartSpanConfig) {
		if cfg.Tags == nil {
			cfg.Tags = map[string]interface{}{}
		}
		cfg.Tags[k] = v
	}
}

// ServiceName sets the given service name on the started span. For example "http.server".
func ServiceName(name string) StartSpanOption {
	return Tag(ext.ServiceName, name)
}

// ResourceName sets the given resource name on the started span. A resource could
// be an SQL query, a URL, an RPC method or something else.
func ResourceName(name string) StartSpanOption {
	return Tag(ext.ResourceName, name)
}

// SpanType sets the given span type on the started span. Some examples in the case of
// the Datadog APM product could be "web", "db" or "cache".
func SpanType(name string) StartSpanOption {
	return Tag(ext.SpanType, name)
}

// WithSpanLinks sets span links on the started span.
func WithSpanLinks(links []ddtrace.SpanLink) StartSpanOption {
	return v2.WithSpanLinks(links)
}

// Measured marks this span to be measured for metrics and stats calculations.
func Measured() StartSpanOption {
	// cache a global instance of this tag: saves one alloc/call
	return v2.Measured()
}

// WithSpanID sets the SpanID on the started span, instead of using a random number.
// If there is no parent Span (eg from ChildOf), then the TraceID will also be set to the
// value given here.
func WithSpanID(id uint64) StartSpanOption {
	return v2.WithSpanID(id)
}

// ChildOf tells StartSpan to use the given span context as a parent for the
// created span.
func ChildOf(ctx ddtrace.SpanContext) StartSpanOption {
	return v2.ChildOf(v2.FromGenericCtx(&internal.SpanContextV1Adapter{Ctx: ctx}))
}

// withContext associates the ctx with the span.
func withContext(ctx context.Context) StartSpanOption {
	return func(cfg *ddtrace.StartSpanConfig) {
		cfg.Context = ctx
	}
}

// StartTime sets a custom time as the start time for the created span. By
// default a span is started using the creation time.
func StartTime(t time.Time) StartSpanOption {
	return v2.StartTime(t)
}

// AnalyticsRate sets a custom analytics rate for a span. It decides the percentage
// of events that will be picked up by the App Analytics product. It's represents a
// float64 between 0 and 1 where 0.5 would represent 50% of events.
func AnalyticsRate(rate float64) StartSpanOption {
	return v2.AnalyticsRate(rate)
}

// FinishOption is a configuration option for FinishSpan. It is aliased in order
// to help godoc group all the functions returning it together. It is considered
// more correct to refer to it as the type as the origin, ddtrace.FinishOption.
type FinishOption = ddtrace.FinishOption

// FinishTime sets the given time as the finishing time for the span. By default,
// the current time is used.
func FinishTime(t time.Time) FinishOption {
	return v2.FinishTime(t)
}

// WithError marks the span as having had an error. It uses the information from
// err to set tags such as the error message, error type and stack trace. It has
// no effect if the error is nil.
func WithError(err error) FinishOption {
	return v2.WithError(err)
}

// NoDebugStack prevents any error presented using the WithError finishing option
// from generating a stack trace. This is useful in situations where errors are frequent
// and performance is critical.
func NoDebugStack() FinishOption {
	return v2.NoDebugStack()
}

// StackFrames limits the number of stack frames included into erroneous spans to n, starting from skip.
func StackFrames(n, skip uint) FinishOption {
	return v2.StackFrames(n, skip)
}

// WithHeaderTags enables the integration to attach HTTP request headers as span tags.
// Warning:
// Using this feature can risk exposing sensitive data such as authorization tokens to Datadog.
// Special headers can not be sub-selected. E.g., an entire Cookie header would be transmitted, without the ability to choose specific Cookies.
func WithHeaderTags(headerAsTags []string) StartOption {
	return v2.WithHeaderTags(headerAsTags)
}

// setHeaderTags sets the global header tags.
// Always resets the global value and returns true.
func setHeaderTags(headerAsTags []string) bool {
	globalconfig.ClearHeaderTags()
	for _, h := range headerAsTags {
		if strings.HasPrefix(h, "x-datadog-") {
			continue
		}
		header, tag := normalizer.HeaderTag(h)
		globalconfig.SetHeaderTag(header, tag)
	}
	return true
}

// UserMonitoringConfig is used to configure what is used to identify a user.
// This configuration can be set by combining one or several UserMonitoringOption with a call to SetUser().
type UserMonitoringConfig struct {
	PropagateID bool
	Email       string
	Name        string
	Role        string
	SessionID   string
	Scope       string
	Metadata    map[string]string
}

// UserMonitoringOption represents a function that can be provided as a parameter to SetUser.
type UserMonitoringOption func(*UserMonitoringConfig)

// WithUserMetadata returns the option setting additional metadata of the authenticated user.
// This can be used multiple times and the given data will be tracked as `usr.{key}=value`.
func WithUserMetadata(key, value string) UserMonitoringOption {
	return func(cfg *UserMonitoringConfig) {
		cfg.Metadata[key] = value
	}
}

// WithUserEmail returns the option setting the email of the authenticated user.
func WithUserEmail(email string) UserMonitoringOption {
	return func(cfg *UserMonitoringConfig) {
		cfg.Email = email
	}
}

// WithUserName returns the option setting the name of the authenticated user.
func WithUserName(name string) UserMonitoringOption {
	return func(cfg *UserMonitoringConfig) {
		cfg.Name = name
	}
}

// WithUserSessionID returns the option setting the session ID of the authenticated user.
func WithUserSessionID(sessionID string) UserMonitoringOption {
	return func(cfg *UserMonitoringConfig) {
		cfg.SessionID = sessionID
	}
}

// WithUserRole returns the option setting the role of the authenticated user.
func WithUserRole(role string) UserMonitoringOption {
	return func(cfg *UserMonitoringConfig) {
		cfg.Role = role
	}
}

// WithUserScope returns the option setting the scope (authorizations) of the authenticated user.
func WithUserScope(scope string) UserMonitoringOption {
	return func(cfg *UserMonitoringConfig) {
		cfg.Scope = scope
	}
}

// WithPropagation returns the option allowing the user id to be propagated through distributed traces.
// The user id is base64 encoded and added to the datadog propagated tags header.
// This option should only be used if you are certain that the user id passed to `SetUser()` does not contain any
// personal identifiable information or any kind of sensitive data, as it will be leaked to other services.
func WithPropagation() UserMonitoringOption {
	return func(cfg *UserMonitoringConfig) {
		cfg.PropagateID = true
	}
}

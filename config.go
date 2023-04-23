package ott

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

/* Config struct for initializing the tracer.
Attributes -->
AppName : This is the application name
Provider : The tracer provider selected.
Endpoint : For a full provisioned provider setup, the url to the distributed backend tcp listener.
SkipExport : SkipExport basically skips exporting of provider and defaults to noop.
Providers : Map of provider constructor.
*/

type Config struct {
	AppName string `json:"appName"`
	Provider string `json:"provider"`
	Endpoint string `json:"endpoint"`
	SkipExport bool `json:"skipExport"`
	Providers map[string]ProviderConstructor `json"-"`
}

type TraceConfig struct {
	TraceProvider trace.TracerProvider
}

// Interface for the TracerProvider and propagator.

type Tracing struct {
	// TracerProvider helps create trace spans.
	TracerProvider trace.TracerProvider
	// Propagator helps propagate trace context across API boundaries as a cross-cutting concern.
	Propagator propagation.TextMapPropagator
}
package ott

import (
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

/*
Initializes the tracing struct with provider tracing provider in the config.
 */

func InitTracing(c Config, appName string) (Tracing, error) {
	var tracing = Tracing{
		Propagator: propagation.TraceContext{},
		TracerProvider: trace.NewNoopTracerProvider(),
	}
	c.AppName = appName
	tracerProvider, err := ConfigureTracerProvider(c)
	if err != nil {
		return Tracing{}, err
	}
	tracing.TracerProvider = tracerProvider
	return tracing, nil
}

// Returns the mux router middlware options to be added to enable auto instrumentation of tracing.

func GetMuxOptions(tr Tracing) []otelmux.Option {
	return []otelmux.Option{
		otelmux.WithPropagators(tr.Propagator),
		otelmux.WithTracerProvider(tr.TracerProvider),
	}
}


// Push the OpenTelemetry transport configuration inside the router client.

func NewTransport(tr Tracing) *otelhttp.Transport {
	return otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithPropagators(tr.Propagator),
		otelhttp.WithTracerProvider(tr.TracerProvider),
		)
}




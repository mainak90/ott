package ott

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.13.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrTracerProviderNotFound    = errors.New("TracerProvider could not be found")
	ErrTracerProviderBuildFailed = errors.New("Failed building TracerProvider")
)

/* DefaultTracerProvider is used when no provider is given.
The Noop tracer provider turns all tracing related operations into noops essentially disabling tracing. */

const DefaultTracerProvider = "noop"

/* ConfigureTracerProvider returns the TracerProvider based on the configuration provided.
It has built-in support for jaeger, zipkin, stdout and noop providers.
A different provider can be used if a constructor for it is provided in the
config. If a provider name is not provided, a noop tracerProvider will be returned. */

func ConfigureTracerProvider(config Config) (trace.TracerProvider, error) {
	if len(config.Provider) == 0 {
		config.Provider = DefaultTracerProvider
	}
	// Handling camelcase of provider.
	config.Provider = strings.ToLower(config.Provider)
	providerConfig := config.Providers[config.Provider]
	if providerConfig == nil {
		providerConfig = providersConfig[config.Provider]
	}
	if providerConfig == nil {
		return nil, fmt.Errorf("%w for provider %s", ErrTracerProviderNotFound, config.Provider)
	}
	provider, err := providerConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTracerProviderBuildFailed, err)
	}
	return provider, nil
}

// ProviderConstructor is useful when client wants to add their own custom TracerProvider.
type ProviderConstructor func(config Config) (trace.TracerProvider, error)

// Created pre-defined immutable map of built-in provider's
var providersConfig = map[string]ProviderConstructor{
	"jaeger": func(cfg Config) (trace.TracerProvider, error) {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
		if err != nil {
			return nil, err
		}
		traceProvider := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp), sdktrace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL, semconv.ServiceNameKey.String(cfg.AppName),
					attribute.String("exporter", cfg.Provider))),
			)
		return traceProvider, nil
	},
	"zipkin": func(cfg Config) (trace.TracerProvider, error) {
		var logger = log.New(os.Stderr, cfg.AppName, log.Ldate|log.Ltime|log.Llongfile)
		exporter, err := zipkin.New(
			cfg.Endpoint,
			zipkin.WithLogger(logger),
		)
		if err != nil {
			return nil, err
		}
		batcher := sdktrace.NewBatchSpanProcessor(exporter)
		traceProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSpanProcessor(batcher),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.AppName),
			)),
		)
		return traceProvider, err
	},
	"stdout": func(cfg Config) (trace.TracerProvider, error) {
		exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to initialize stdouttrace exporter: %w", err)
		}
		batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exp)
		traceProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(batchSpanProcessor),
		)
		return traceProvider, nil
	},
	"noop": func(config Config) (trace.TracerProvider, error) {
		return trace.NewNoopTracerProvider(), nil
	},
}
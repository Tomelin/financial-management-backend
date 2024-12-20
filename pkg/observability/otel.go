package observability

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oTrace "go.opentelemetry.io/otel/trace"
)

type OtelProvider struct {
	Name     string            `json:"name" mapstructure:"name"`
	Endpoint string            `json:"endpoint,omitempty" mapstructure:"endpoint"`
	Headers  map[string]string `json:"headers,omitempty" mapstructure:"headers"`
	Insecure bool              `json:"insecure,omitempty" mapstructure:"insecure"`
}

type OtelConfig struct {
	ServiceNameKey         string         `json:"service_name_key" mapstructure:"service_name_key"`
	ServiceVersionKey      string         `json:"service_version_key" mapstructure:"service_version_key"`
	ServiceEnvironmentKey  string         `json:"service_environment_key" mapstructure:"service_environment_key"`
	ServiceEndpoint        string         `json:"service_endpoint" mapstructure:"service_endpoint"`
	Providers              []OtelProvider `json:"providers" mapstructure:"providers"`
	TraceparentHeaderName  string         `json:"traceparent_header_name" mapstructure:"traceparent_header_name"`
	TracestateHeaderName   string         `json:"tracestate_header_name" mapstructure:"tracestate_header_name"`
	TracecontextHeaderName string         `json:"tracecontext_header_name" mapstructure:"tracecontext_header_name"`
	TraceidHeaderName      string         `json:"traceid_header_name" mapstructure:"traceid_header_name"`
	SpanidHeaderName       string         `json:"spanid_header_name" mapstructure:"spanid_header_name"`
	SamplingPolicy         string         `json:"sampling_policy" mapstructure:"sampling_policy"`
	SpanContextHeaderName  string         `json:"span_context_header_name" mapstructure:"span_context_header_name"`
}

type Tracer struct {
	Trace oTrace.Tracer
}

func InicializeTracer(fields any) (*Tracer, func(), error) {

	config, err := parseConfig(fields)
	if err != nil {
		return nil, nil, err
	}

	var spanProcessors []trace.TracerProviderOption

	exporterInitializers := map[string]func(OtelProvider) (trace.SpanExporter, error){
		"stdout": func(provider OtelProvider) (trace.SpanExporter, error) {
			return stdouttrace.New(stdouttrace.WithPrettyPrint())
		},
		"jaeger": func(provider OtelProvider) (trace.SpanExporter, error) {
			return otlptracehttp.New(context.Background(),
				otlptracehttp.WithEndpoint(provider.Endpoint),
				otlptracehttp.WithHeaders(provider.Headers),
				otlptracehttp.WithInsecure(),
			)
		},
	}

	for _, provider := range config.Providers {
		if initializer, ok := exporterInitializers[provider.Name]; ok {
			exporter, err := initializer(provider)
			if err != nil {
				return nil, nil, err
			}
			spanProcessors = append(spanProcessors, trace.WithBatcher(exporter))
		}
	}

	// Configura o provedor de rastreamento com os exportadores habilitados
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceNameKey),
		semconv.ServiceVersionKey.String(config.ServiceVersionKey),
		semconv.DeploymentEnvironmentKey.String("staging"),
	)

	traceProvider := trace.NewTracerProvider(
		append(spanProcessors, trace.WithResource(resources))...,
	)

	otel.SetTracerProvider(traceProvider)

	cleanup := func() {
		traceProvider.Shutdown(context.Background())
	}
	t := otel.Tracer(config.ServiceNameKey)
	return &Tracer{Trace: t}, cleanup, nil
}

func parseConfig(fields any) (*OtelConfig, error) {

	b, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	var config OtelConfig
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

package shared

import "go.opentelemetry.io/otel/trace"

type TemplateData struct {
	RequestNameOTEL string
	OTELTracer      trace.Tracer
}

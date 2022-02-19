package w3c

// from: https://github.com/open-telemetry/opentelemetry-go/blob/v1.4.0/propagation/trace_context.go#L26-L31
const (
	SupportedVersion  = 0
	MaxVersion        = 254
	TraceparentHeader = "traceparent"
	TracestateHeader  = "tracestate"
)

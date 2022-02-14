package w3c

import (
	"fmt"
	"go.opentelemetry.io/otel/trace"
)

// TraceParent implements the W3C trace-context standard: https://w3c.github.io/trace-context/
type TraceParent struct {
	ContextVersion string
	TraceID        trace.TraceID
	ParentID       trace.SpanID
	TraceFlags     trace.TraceFlags
}

// NewTraceParentFromSpanContext creates a new TraceParent from the given trace.SpanContext
func NewTraceParentFromSpanContext(ctx trace.SpanContext) TraceParent {
	return TraceParent{
		ContextVersion: "00",
		TraceID:        ctx.TraceID(),
		ParentID:       ctx.SpanID(),
		TraceFlags:     ctx.TraceFlags(),
	}
}

// String implements the Stringer interface for TraceParent
func (t TraceParent) String() string {
	return fmt.Sprintf("%s-%s-%s-%s", t.ContextVersion, t.TraceID.String(), t.ParentID.String(), t.TraceFlags.String())
}

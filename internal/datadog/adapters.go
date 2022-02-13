package datadog

import (
	"encoding/hex"
	"strconv"
)

// OpenTelemetry TraceId and SpanId properties differ from Datadog conventions. Therefore, it’s necessary to translate
// TraceId and SpanId from their OpenTelemetry formats (a 128bit unsigned int and 64bit unsigned int represented as a
// 32-hex-character and 16-hex-character lowercase string, respectively) into their Datadog Formats(a 64bit unsigned
// int).
// - https: //docs.datadoghq.com/tracing/connect_logs_and_traces/opentelemetry/
//
// The following adapter functions are copied from:
// - https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/v0.44.0/exporter/datadogexporter/translate_traces.go#L491-L508

// DecodeAPMTraceID maps an OpenTelemetry hexidecimal TraceID into a Datadog uint64 TraceID
func DecodeAPMTraceID(rawID [16]byte) uint64 {
	return DecodeAPMId(hex.EncodeToString(rawID[:]))
}

// DecodeAPMSpanID maps an OpenTelemetry hexidecimal SpanID into a Datadog uint64 SpanID
func DecodeAPMSpanID(rawID [8]byte) uint64 {
	return DecodeAPMId(hex.EncodeToString(rawID[:]))
}

// DecodeAPMId maps an OpenTelemetry string into a Datadog formatted uint64
func DecodeAPMId(id string) uint64 {
	if len(id) > 16 {
		id = id[len(id)-16:]
	}
	val, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return 0
	}
	return val
}

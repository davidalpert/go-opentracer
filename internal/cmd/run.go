package cmd

import (
	"context"
	"fmt"
	"github.com/davidalpert/gopentracer/internal/datadog"
	"github.com/davidalpert/gopentracer/internal/types"
	"github.com/davidalpert/gopentracer/internal/utils"
	"github.com/davidalpert/gopentracer/internal/version"
	"github.com/davidalpert/gopentracer/internal/w3c"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// RunOptions is a struct to support version command
type RunOptions struct {
	utils.PrinterOptions
	Command               string
	DeploymentEnvironment string
	SpanName              string
	SpanTagsRaw           []string
	TraceOLTPHttpEndpoint string
	TraceLogFile          string
	SpanDelay             time.Duration
	VersionSummary        version.SummaryStruct
}

// NewRunOptions returns initialized RunOptions
func NewRunOptions() *RunOptions {
	return &RunOptions{
		VersionSummary: version.Summary,
	}
}

// NewCmdRun creates the version command
func NewCmdRun() *cobra.Command {
	o := NewRunOptions()
	var cmd = &cobra.Command{
		Use:   "run <cmd>",
		Short: "runs a command inside an open trace and span",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err

			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	o.AddPrinterFlags(cmd)
	cmd.Flags().StringVarP(&o.DeploymentEnvironment, "deployment-environment", "e", "prd", "deployment environment")
	cmd.Flags().StringVar(&o.TraceOLTPHttpEndpoint, "trace-http-endpoint", "", "sent traces over http to this endpoint")
	cmd.Flags().StringVar(&o.TraceLogFile, "trace-log-file", "", "log traces to this file")
	cmd.Flags().StringSliceVar(&o.SpanTagsRaw, "tag", make([]string, 0), "tags in the format key:val[:type]")
	cmd.Flags().DurationVar(&o.SpanDelay, "span-delay", 100*time.Millisecond, "how long to wait after the command completes before completing the span (golang time.Duration)")
	cmd.Flags().StringVar(&o.SpanName, "span-name", "Run", "name for this span")
	return cmd
}

// Complete completes the RunOptions
func (o *RunOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Command = args[0]
	return nil
}

// Validate validates the RunOptions
func (o *RunOptions) Validate() error {
	if o.Command == "" {
		return fmt.Errorf("command is required")
	}
	if o.SpanName == "" {
		return fmt.Errorf("span-name is required")
	}
	return o.PrinterOptions.Validate()
}

// newConsoleExporter returns a console exporter.
func newConsoleExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func newFileExporter(filename string) (*sdktrace.SpanExporter, func(), error) {
	cleanupFN := func() {}
	if filename == "" {
		return nil, cleanupFN, fmt.Errorf("cannot export to an empty filename")
	}
	// Write telemetry data to a file.
	f, err := os.Create(filename)
	if err != nil {
		return nil, cleanupFN, err
	}
	cleanupFN = func() { f.Close() }

	exp, err := newConsoleExporter(f)
	if err != nil {
		return nil, cleanupFN, err
	}

	return &exp, cleanupFN, nil
}

// newTracerResource returns a resource describing this application.
func (o *RunOptions) newTracerResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(o.VersionSummary.AppName),
			semconv.ServiceVersionKey.String(o.VersionSummary.Version),
			semconv.DeploymentEnvironmentKey.String(o.DeploymentEnvironment),
		),
	)
	return r
}

// Run executes the command
func (o *RunOptions) Run() error {
	traceProviderOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(o.newTracerResource()),
	}

	if o.TraceLogFile != "" {
		exp, cleanupFN, err := newFileExporter(o.TraceLogFile)
		if err != nil {
			return err
		}
		defer cleanupFN()
		traceProviderOptions = append(traceProviderOptions, sdktrace.WithBatcher(*exp))
	}

	if o.TraceOLTPHttpEndpoint != "" {
		exp, err := otlptracehttp.New(context.TODO(),
			buildHttpTraceExporterSpanOptionsForEndpoint(o.TraceOLTPHttpEndpoint)...,
		)
		if err != nil {
			return err
		}
		traceProviderOptions = append(traceProviderOptions, sdktrace.WithBatcher(exp))
	}

	tp := sdktrace.NewTracerProvider(traceProviderOptions...)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	otel.SetTracerProvider(tp)

	ctx, span := otel.Tracer(o.VersionSummary.AppName,
		trace.WithInstrumentationVersion(o.VersionSummary.Version),
	).Start(context.TODO(), o.SpanName)
	defer span.End()
	cmdCtx := trace.ContextWithSpan(context.TODO(), span)

	for _, s := range o.SpanTagsRaw {
		if a, err := rawTagToTypedAttribute(cmdCtx, s); err != nil {
			return err
		} else {
			span.SetAttributes(a)
		}
	}

	cmdText := injectTraceAndSpanID(cmdCtx, o.Command)
	//cmdText = os.ExpandEnv(cmdText) // TODO: review as this may create unpredictable behavior
	cmdParts := strings.Split(cmdText, " ")
	c := exec.CommandContext(ctx, cmdParts[0], cmdParts[1:]...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Env = make([]string, len(os.Environ()))
	for i, e := range os.Environ() {
		c.Env[i] = e
	}
	c.Env = appendTraceAndSpanIDToEnv(ctx, c.Env)

	err := c.Run()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	time.Sleep(o.SpanDelay)

	return err
}

func buildHttpTraceExporterSpanOptionsForEndpoint(endpoint string) []otlptracehttp.Option {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}

	if !strings.HasPrefix(endpoint, "https://") {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	return opts
}

func injectTraceAndSpanID(ctx context.Context, s string) string {
	// open telemetry always returns a span; if the given ctx doesn't have one
	// then trace.SpanFromContext returns a noopspan which implements trace.Span
	// but has no data and no functionality; thus it is safe to use as a trace.Span
	spanCtx := trace.SpanFromContext(ctx).SpanContext()

	s = strings.Replace(s, "$TRACE_ID", spanCtx.TraceID().String(), -1)
	s = strings.Replace(s, "${TRACE_ID}", spanCtx.TraceID().String(), -1)

	s = strings.Replace(s, "$SPAN_ID", spanCtx.SpanID().String(), -1)
	s = strings.Replace(s, "${SPAN_ID}", spanCtx.SpanID().String(), -1)

	s = strings.Replace(s, "$PARENT_ID", spanCtx.SpanID().String(), -1)
	s = strings.Replace(s, "${PARENT_ID}", spanCtx.SpanID().String(), -1)

	ddTraceID := fmt.Sprintf("%d", datadog.DecodeAPMTraceID(spanCtx.TraceID()))
	s = strings.Replace(s, "$DD_TRACE_ID", ddTraceID, -1)
	s = strings.Replace(s, "${DD_TRACE_ID}", ddTraceID, -1)

	ddSpanID := fmt.Sprintf("%d", datadog.DecodeAPMSpanID(spanCtx.SpanID()))
	s = strings.Replace(s, "$DD_SPAN_ID", ddSpanID, -1)
	s = strings.Replace(s, "${DD_SPAN_ID}", ddSpanID, -1)

	s = strings.Replace(s, "$DD_PARENT_ID", ddSpanID, -1)
	s = strings.Replace(s, "${DD_PARENT_ID}", ddSpanID, -1)

	traceparentValue := w3c.NewTraceParentFromSpanContext(spanCtx).String()
	s = strings.Replace(s, "$W3CTRACEPARENT", traceparentValue, -1)
	s = strings.Replace(s, "${W3CTRACEPARENT}", traceparentValue, -1)

	return s
}

func appendTraceAndSpanIDToEnv(ctx context.Context, ss []string) []string {
	ss = append(ss, injectTraceAndSpanID(ctx, "TRACE_ID=$TRACE_ID"))
	ss = append(ss, injectTraceAndSpanID(ctx, "SPAN_ID=$SPAN_ID"))
	ss = append(ss, injectTraceAndSpanID(ctx, "DD_TRACE_ID=$DD_TRACE_ID"))
	ss = append(ss, injectTraceAndSpanID(ctx, "DD_SPAN_ID=$DD_SPAN_ID"))
	ss = append(ss, injectTraceAndSpanID(ctx, "W3CTRACEPARENT=$W3CTRACEPARENT"))
	ss = append(ss, fmt.Sprintf("GOPENTRACER_VERSION=%s", version.Summary.Version))
	return ss
}

func rawTagToTypedAttribute(ctx context.Context, s string) (attribute.KeyValue, error) {
	parts := strings.Split(s, ":")
	if len(parts) == 2 {
		// value without type, default to string
		parts = append(parts, "string")
	}
	if len(parts) < 3 {
		return attribute.KeyValue{}, fmt.Errorf("must specify key:value (or optionally key:value:type): '%s'", s)
	}

	key := parts[0]
	val := injectTraceAndSpanID(ctx, parts[1])
	valType := parts[2]

	var attrType types.OpenTelemetryAttributeType
	err := attrType.UnmarshalJSON([]byte("\"" + valType + "\""))
	if err != nil {
		return attribute.KeyValue{}, err
	}

	switch attrType {
	case types.StringAttribute:
		return attribute.String(key, val), nil
	case types.BoolAttribute:
		if v, err := strconv.ParseBool(val); err != nil {
			return attribute.String(key, val), nil
		} else {
			return attribute.Bool(key, v), nil
		}
	case types.IntAttribute, types.Int32Attribute:
		if v, err := strconv.ParseInt(val, 10, 32); err != nil {
			return attribute.String(key, val), nil
		} else {
			return attribute.Int(key, int(v)), nil
		}
	case types.Int64Attribute:
		if v, err := strconv.ParseInt(val, 10, 64); err != nil {
			return attribute.String(key, val), nil
		} else {
			return attribute.Int64(key, v), nil
		}
	default:
		panic("should never get here")
	}
}

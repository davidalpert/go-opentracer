package cmd

import (
	"context"
	"fmt"
	"github.com/davidalpert/gopentracer/internal/types"
	"github.com/davidalpert/gopentracer/internal/utils"
	"github.com/davidalpert/gopentracer/internal/version"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"
	"os/exec"
	"strings"
)

// RunOptions is a struct to support version command
type RunOptions struct {
	utils.PrinterOptions
	Args                  []string
	Path                  string
	DeploymentEnvironment string
	SpanTagAttributes     []attribute.KeyValue
	SpanTagsRaw           []string
	TraceOLTPHttpEndpoint string
	TraceLogFile          string
	VersionSummary        version.SummaryStruct
}

// NewRunOptions returns initialized RunOptions
func NewRunOptions() *RunOptions {
	return &RunOptions{
		Args:              make([]string, 0),
		SpanTagAttributes: make([]attribute.KeyValue, 0),
		VersionSummary:    version.Summary,
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
	return cmd
}

// Complete completes the RunOptions
func (o *RunOptions) Complete(cmd *cobra.Command, args []string) error {
	commandParts := strings.Split(args[0], " ")
	o.Path = commandParts[0]
	if len(commandParts) > 1 {
		o.Args = commandParts[1:]
	}

	// TODO: consider move creation of the span and configuration of the tracer here?

	return nil
}

func rawTagToKeyValue(spanCtx trace.SpanContext, s string) (attribute.KeyValue, error) {
	parts := strings.Split(s, ":")
	if len(parts) == 2 {
		// value without type, default to string
		parts = append(parts, "string")
	}
	if len(parts) < 3 {
		return attribute.KeyValue{}, fmt.Errorf("must specify key:value (or optionally key:value:type): '%s'", s)
	}

	key := parts[0]
	val := injectTraceAndSpanID(spanCtx, parts[1])
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

// Validate validates the RunOptions
func (o *RunOptions) Validate() error {
	if o.Path == "" {
		return fmt.Errorf("command is required")
	}
	return o.PrinterOptions.Validate()
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func (o *RunOptions) newResource() *resource.Resource {
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
	l := log.New(os.Stdout, "", 0)

	traceProviderOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(o.newResource()),
	}

	if o.TraceLogFile != "" {
		// Write telemetry data to a file.
		f, err := os.Create(o.TraceLogFile)
		if err != nil {
			l.Fatal(err)
		}
		defer f.Close()

		exp, err := newExporter(f)
		if err != nil {
			l.Fatal(err)
		}
		traceProviderOptions = append(traceProviderOptions, sdktrace.WithBatcher(exp))
	}

	if o.TraceOLTPHttpEndpoint != "" {
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(o.TraceOLTPHttpEndpoint),
		}
		if !strings.HasPrefix(o.TraceOLTPHttpEndpoint, "https://") {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		exp, err := otlptracehttp.New(context.TODO(), opts...)
		if err != nil {
			l.Fatal(err)
		}

		traceProviderOptions = append(traceProviderOptions, sdktrace.WithBatcher(exp))
	}

	tp := sdktrace.NewTracerProvider(traceProviderOptions...)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	ctx, span := otel.Tracer(o.VersionSummary.AppName,
		trace.WithInstrumentationVersion(o.VersionSummary.Version),
	).Start(context.TODO(), "Run",
		trace.WithAttributes(o.SpanTagAttributes...),
	)
	defer span.End()
	spanCtx := span.SpanContext()

	spanAttrs := make([]attribute.KeyValue, len(o.SpanTagsRaw))
	for i, s := range o.SpanTagsRaw {
		if a, err := rawTagToKeyValue(span.SpanContext(), s); err != nil {
			return err
		} else {
			spanAttrs[i] = a
		}
	}

	args := make([]string, len(o.Args))
	for i, s := range o.Args {
		args[i] = injectTraceAndSpanID(spanCtx, s)
	}

	c := exec.CommandContext(ctx, o.Path, args...)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "TRACE_ID:$TRACE_ID"))
	c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "SPAN_ID:$SPAN_ID"))
	c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "DD_TRACE_ID:$DD_TRACE_ID"))
	c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "DD_SPAN_ID:$DD_SPAN_ID"))

	fmt.Printf("path: %s\nargs: %#v\nenv: %#v\ntags: %#v\n", c.Path, c.Args, c.Env, spanAttrs)

	err := c.Run()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

const emptyTraceID = "00000000000000000000000000000000"
const emptySpanID = "0000000000000000"

func hexToUint64(h string) uint64 {
	u, _ := strconv.ParseUint(h, 16, 64)
	return u
}

func injectTraceAndSpanID(spanCtx trace.SpanContext, s string) string {
	s = strings.Replace(s, "$TRACE_ID", spanCtx.TraceID().String(), -1)
	s = strings.Replace(s, "${TRACE_ID}", spanCtx.TraceID().String(), -1)

	s = strings.Replace(s, "$SPAN_ID", spanCtx.SpanID().String(), -1)
	s = strings.Replace(s, "${SPAN_ID}", spanCtx.SpanID().String(), -1)

	ddTraceID := fmt.Sprintf("%d", hexToUint64(spanCtx.TraceID().String()))
	s = strings.Replace(s, "$DD_TRACE_ID", ddTraceID, -1)
	s = strings.Replace(s, "${DD_TRACE_ID}", ddTraceID, -1)

	ddSpanID := fmt.Sprintf("%d", hexToUint64(spanCtx.SpanID().String()))
	s = strings.Replace(s, "$DD_SPAN_ID", ddSpanID, -1)
	s = strings.Replace(s, "${DD_SPAN_ID}", ddSpanID, -1)

	s = strings.Replace(s, "$PARENT_TRACE_ID", emptyTraceID, -1)
	s = strings.Replace(s, "${PARENT_TRACE_ID}", emptyTraceID, -1)

	s = strings.Replace(s, "$PARENT_SPAN_ID", emptySpanID, -1)
	s = strings.Replace(s, "${PARENT_SPAN_ID}", emptySpanID, -1)
	return s
}

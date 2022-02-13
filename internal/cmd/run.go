package cmd

import (
	"context"
	"fmt"
	"github.com/davidalpert/gopentracer/internal/utils"
	"github.com/davidalpert/gopentracer/internal/version"
	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"strconv"

	"go.opentelemetry.io/otel/trace"
	"os"
	"os/exec"
	"strings"
)

// RunOptions is a struct to support version command
type RunOptions struct {
	utils.PrinterOptions
	Command               string
	DeploymentEnvironment string
	SpanTagsRaw           []string
	//TraceOLTPHttpEndpoint string
	//TraceLogFile          string
	VersionSummary version.SummaryStruct
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
		Short: "runs a command inside a datadog trace and span",
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
	cmd.Flags().StringSliceVar(&o.SpanTagsRaw, "tag", make([]string, 0), "tags in the format key:val[:type]")
	//cmd.Flags().StringVar(&o.TraceOLTPHttpEndpoint, "trace-http-endpoint", "", "sent traces over http to this endpoint")
	//cmd.Flags().StringVar(&o.TraceLogFile, "trace-log-file", "", "log traces to this file")
	return cmd
}

// Complete completes the RunOptions
func (o *RunOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Command = args[0]
	return nil
}

func formatCommand(ctx context.Context, cmdText string) *exec.Cmd {
	cmdText = injectTraceAndSpanID2(ctx, cmdText)
	cmdParts := strings.Split(cmdText, " ")
	cmd := exec.CommandContext(ctx, cmdParts[0], cmdParts[1:]...)
	return cmd
}

func injectTraceAndSpanID2(ctx context.Context, s string) string {
	s = injectDatadogSpanAndTraceID(ctx, s)
	return s
}

func injectDatadogSpanAndTraceID(ctx context.Context, s string) string {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		spanCtx := span.Context()
		ddTraceID := fmt.Sprintf("%d", spanCtx.TraceID())
		ddSpanID := fmt.Sprintf("%d", spanCtx.SpanID())

		s = strings.Replace(s, "$DD_TRACE_ID", ddTraceID, -1)
		s = strings.Replace(s, "$DD_SPAN_ID", ddSpanID, -1)
	}
	return s
}

func appendTraceAndSpanIDToEnv(ctx context.Context, ss []string) []string {
	ss = append(ss, injectTraceAndSpanID2(ctx, "DD_TRACE_ID:$DD_TRACE_ID"))
	ss = append(ss, injectTraceAndSpanID2(ctx, "DD_SPAN_ID:$DD_SPAN_ID"))
	return ss
}

//func rawTagToKeyValue(spanCtx trace.SpanContext, s string) (attribute.KeyValue, error) {
//	parts := strings.Split(s, ":")
//	if len(parts) == 2 {
//		// value without type, default to string
//		parts = append(parts, "string")
//	}
//	if len(parts) < 3 {
//		return attribute.KeyValue{}, fmt.Errorf("must specify key:value (or optionally key:value:type): '%s'", s)
//	}
//
//	key := parts[0]
//	val := injectTraceAndSpanID(spanCtx, parts[1])
//	valType := parts[2]
//
//	var attrType types.OpenTelemetryAttributeType
//	err := attrType.UnmarshalJSON([]byte("\"" + valType + "\""))
//	if err != nil {
//		return attribute.KeyValue{}, err
//	}
//
//	switch attrType {
//	case types.StringAttribute:
//		return attribute.String(key, val), nil
//	case types.BoolAttribute:
//		if v, err := strconv.ParseBool(val); err != nil {
//			return attribute.String(key, val), nil
//		} else {
//			return attribute.Bool(key, v), nil
//		}
//	case types.IntAttribute, types.Int32Attribute:
//		if v, err := strconv.ParseInt(val, 10, 32); err != nil {
//			return attribute.String(key, val), nil
//		} else {
//			return attribute.Int(key, int(v)), nil
//		}
//	case types.Int64Attribute:
//		if v, err := strconv.ParseInt(val, 10, 64); err != nil {
//			return attribute.String(key, val), nil
//		} else {
//			return attribute.Int64(key, v), nil
//		}
//	default:
//		panic("should never get here")
//	}
//}

// Validate validates the RunOptions
func (o *RunOptions) Validate() error {
	if o.Command == "" {
		return fmt.Errorf("command is required")
	}
	return o.PrinterOptions.Validate()
}

//// newExporter returns a console exporter.
//func newExporter(w io.Writer) (sdktrace.SpanExporter, error) {
//	return stdouttrace.New(
//		stdouttrace.WithWriter(w),
//		// Use human readable output.
//		stdouttrace.WithPrettyPrint(),
//		// Do not print timestamps for the demo.
//		stdouttrace.WithoutTimestamps(),
//	)
//}

//// newResource returns a resource describing this application.
//func (o *RunOptions) newResource() *resource.Resource {
//	r, _ := resource.Merge(
//		resource.Default(),
//		resource.NewWithAttributes(
//			semconv.SchemaURL,
//			semconv.ServiceNameKey.String(o.VersionSummary.AppName),
//			semconv.ServiceVersionKey.String(o.VersionSummary.Version),
//			semconv.DeploymentEnvironmentKey.String(o.DeploymentEnvironment),
//		),
//	)
//	return r
//}

// Run executes the command
func (o *RunOptions) Run() error {
	//l := log.New(os.Stdout, "", 0)

	tracer.Start(
		tracer.WithEnv(o.DeploymentEnvironment),
		tracer.WithService(o.VersionSummary.AppName),
		tracer.WithServiceVersion(o.VersionSummary.Version),
	)

	// When the tracer is stopped, it will flush everything it has to the Datadog Agent before quitting.
	// Make sure this line stays in your main function.
	defer tracer.Stop()

	span := tracer.StartSpan("run")
	spanFinishOpts := make([]tracer.FinishOption, 0)
	defer span.Finish(spanFinishOpts...)

	for _, rt := range o.SpanTagsRaw {
		parts := strings.Split(rt, ":")
		if len(parts) == 2 {
			span.SetTag(parts[0], parts[1])
		}
	}

	cmdCtx := tracer.ContextWithSpan(context.TODO(), span)
	cmd := formatCommand(cmdCtx, o.Command)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = appendTraceAndSpanIDToEnv(cmdCtx, cmd.Env)

	err := cmd.Run()
	if err != nil {
		spanFinishOpts = append(spanFinishOpts, tracer.WithError(err))
	}

	return err
	//
	//traceProviderOptions := []sdktrace.TracerProviderOption{
	//	sdktrace.WithResource(o.newResource()),
	//}
	//
	//if o.TraceLogFile != "" {
	//	// Write telemetry data to a file.
	//	f, err := os.Create(o.TraceLogFile)
	//	if err != nil {
	//		l.Fatal(err)
	//	}
	//	defer f.Close()
	//
	//	exp, err := newExporter(f)
	//	if err != nil {
	//		l.Fatal(err)
	//	}
	//	traceProviderOptions = append(traceProviderOptions, sdktrace.WithBatcher(exp))
	//}
	//
	//if o.TraceOLTPHttpEndpoint != "" {
	//	opts := []otlptracehttp.Option{
	//		otlptracehttp.WithEndpoint(o.TraceOLTPHttpEndpoint),
	//	}
	//	if !strings.HasPrefix(o.TraceOLTPHttpEndpoint, "https://") {
	//		opts = append(opts, otlptracehttp.WithInsecure())
	//	}
	//
	//	exp, err := otlptracehttp.New(context.TODO(), opts...)
	//	if err != nil {
	//		l.Fatal(err)
	//	}
	//
	//	traceProviderOptions = append(traceProviderOptions, sdktrace.WithBatcher(exp))
	//}
	//
	//tp := sdktrace.NewTracerProvider(traceProviderOptions...)
	//defer func() {
	//	if err := tp.Shutdown(context.Background()); err != nil {
	//		l.Fatal(err)
	//	}
	//}()
	//otel.SetTracerProvider(tp)
	//
	//ctx, span := otel.Tracer(o.VersionSummary.AppName,
	//	trace.WithInstrumentationVersion(o.VersionSummary.Version),
	//).Start(context.TODO(), "Run",
	//	trace.WithAttributes(o.SpanTagAttributes...),
	//)
	//defer span.End()
	//spanCtx := span.SpanContext()
	//
	//spanAttrs := make([]attribute.KeyValue, len(o.SpanTagsRaw))
	//for i, s := range o.SpanTagsRaw {
	//	if a, err := rawTagToKeyValue(span.SpanContext(), s); err != nil {
	//		return err
	//	} else {
	//		spanAttrs[i] = a
	//	}
	//}
	//
	//args := make([]string, len(o.Args))
	//for i, s := range o.Args {
	//	args[i] = injectTraceAndSpanID(spanCtx, s)
	//}
	//
	//c := exec.CommandContext(ctx, o.Path, args...)
	//c.Stdout = os.Stdout
	//c.Stdin = os.Stdin
	//c.Stderr = os.Stderr
	//c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "TRACE_ID:$TRACE_ID"))
	//c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "SPAN_ID:$SPAN_ID"))
	//c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "DD_TRACE_ID:$DD_TRACE_ID"))
	//c.Env = append(c.Env, injectTraceAndSpanID(spanCtx, "DD_SPAN_ID:$DD_SPAN_ID"))
	//
	//fmt.Printf("path: %s\nargs: %#v\nenv: %#v\ntags: %#v\n", c.Path, c.Args, c.Env, spanAttrs)
	//
	//err := c.Run()
	//if err != nil {
	//	span.RecordError(err)
	//	span.SetStatus(codes.Error, err.Error())
	//}
	//return err
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

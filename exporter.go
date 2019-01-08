package iowriter

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ----------- Options -----------

// WithWriter allows you to set a writer other than the default of os.Stdout
func WithWriter(w io.Writer) func(*exporter) error {
	return func(e *exporter) error {
		e.writer = zapcore.AddSync(w)
		return nil
	}
}

// WithKeyValues is a simple way to add additional "baggage" to all trace and
// metrics logging.  A common key value pair to add would be `serviceName`
// and the name of the service
func WithKeyValues(m map[string]interface{}) func(*exporter) error {
	return func(e *exporter) error {
		e.kv = m
		return nil
	}
}

// WithMessage allows you to override the default "msg":"opencensus" log output
func WithMessage(msg string) func(*exporter) error {
	return func(e *exporter) error {
		e.msg = msg
		return nil
	}
}

// WithDataKey allows you to override the default "opencensus":{...} key that
// contains the "SpanData" and "ViewData" opencensus data
func WithDataKey(key string) func(*exporter) error {
	return func(e *exporter) error {
		e.dataKey = key
		return nil
	}
}

// ------------ New ---------------

// New creates a new exporter that can used with trace.RegisterExporter()
// and view.RegisterExporter().  Use the iowriter.With* options to customize
// additional values
func New(options ...func(*exporter) error) (Exporter, error) {

	e := exporter{
		msg:     "opencensus",
		dataKey: "opencensus",
	}

	var err error
	for _, option := range options {
		err = option(&e)
		if err != nil {
			return nil, errors.Wrap(err, "Could not apply options")
		}
	}

	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	if e.writer == nil {
		e.writer = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		e.writer,
		zapcore.InfoLevel,
	)

	e.logger = zap.New(core)

	if len(e.kv) > 0 {
		fields := []zapcore.Field{}
		for key, value := range e.kv {
			fields = append(fields, zap.Any(key, value))
		}
		e.logger = e.logger.With(fields...)
	}

	return &e, nil
}

// --------------------------------

type exporter struct {
	logger *zap.Logger
	writer zapcore.WriteSyncer
	kv     map[string]interface{}

	msg     string
	dataKey string
}

// Compile time assertion that the exporter implements trace.Exporter, view.Exporter
var _ trace.Exporter = (*exporter)(nil)
var _ view.Exporter = (*exporter)(nil)

func (e *exporter) ExportSpan(sd *trace.SpanData) {

	oc := OpenCensus{
		SpanData: newSpanData(sd),
	}

	defer e.logger.Sync()
	e.logger.Info(e.msg,
		zap.Any(e.dataKey, oc),
	)
}

func (e *exporter) ExportView(vd *view.Data) {
	oc := OpenCensus{
		ViewData: newViewData(vd),
	}

	defer e.logger.Sync()
	e.logger.Info(e.msg,
		zap.Any(e.dataKey, oc),
	)
}

func newSpanData(sd *trace.SpanData) *SpanData {
	if sd == nil {
		return &SpanData{}
	}

	return &SpanData{
		SpanData: sd,
		// Name:         sd.Name,
		TraceID:      sd.TraceID.String(),
		SpanID:       sd.SpanID.String(),
		ParentSpanID: sd.ParentSpanID.String(),

		Duration: sd.EndTime.Sub(sd.StartTime),
	}
}

func newViewData(vd *view.Data) *ViewData {
	return &ViewData{Data: vd}
}

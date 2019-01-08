package iowriter

import (
	"time"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type Exporter interface {
	ExportSpan(sd *trace.SpanData)
	ExportView(vd *view.Data)
}

type OpenCensus struct {
	SpanData *SpanData `json:",omitempty"`
	ViewData *ViewData `json:",omitempty"`
}

type SpanData struct {
	*trace.SpanData

	// Override a few of the values for cleaner marshalling
	TraceID      string
	SpanID       string
	ParentSpanID string

	// https://github.com/ExpediaDotCom/haystack-idl/blob/master/proto/span.proto#L25
	Duration time.Duration
}

type ViewData struct {
	*view.Data
}

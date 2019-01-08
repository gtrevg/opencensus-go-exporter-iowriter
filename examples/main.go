// Package main contains a modified version of the code found at:
//  * https://opencensus.io/exporters/custom-exporter/go/trace/
//  * https://opencensus.io/exporters/custom-exporter/go/metrics/
//
// This example shows
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	iowriter "github.com/gtrevg/opencensus-go-exporter-iowriter"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

func main() {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Please remember to register your exporter
	// so that it can receive exported spanData.

	e, err := iowriter.New(
		iowriter.WithWriter(os.Stdout),
		iowriter.WithKeyValues(map[string]interface{}{
			"k1": "v1",
			"k2": 2,
			"k3": true,
		}),
	)

	if err != nil {
		panic(err)
	}

	// Register trace exporter
	trace.RegisterExporter(e)

	// We need to have registered at least one view
	if err := view.Register(loopCountView); err != nil {
		log.Fatalf("Failed to register loopCountView: %v", err)
	}

	// Register view exporter
	view.RegisterExporter(e)
	view.SetReportingPeriod(100 * time.Millisecond)

	ctx, _ := tag.New(context.Background(), tag.Upsert(keyMethod, "main"))

	for i := int64(0); i < 5; i++ {
		_, span := trace.StartSpan(ctx, fmt.Sprintf("sample-%d", i))
		span.Annotate([]trace.Attribute{trace.Int64Attribute("invocations", 1)}, "Invoked it")
		stats.Record(ctx, mLoops.M(i))
		span.End()
		<-time.After(10 * time.Millisecond)
	}
	<-time.After(500 * time.Millisecond)
}

// The measure and view to be used for demo purposes
var keyMethod, _ = tag.NewKey("method")
var mLoops = stats.Int64("demo/loop_iterations", "The number of loop iterations", "1")
var loopCountView = &view.View{
	Measure: mLoops, Name: "demo/loop_iterations",
	Description: "Number of loop iterations",
	Aggregation: view.Count(),
	TagKeys:     []tag.Key{keyMethod},
}

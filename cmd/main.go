package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var zipkinAddress = flag.String("zipkin-address", "http://localhost:9411/api/v1/spans", "Zipkin spans API url.")

func main() {
	flag.Parse()

	// Create our HTTP collector.
	collector, err := zipkin.NewHTTPCollector(*zipkinAddress)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:8080", "otracegen")

	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(false),
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// Explicitly set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(tracer)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		wireContext, _ := tracer.Extract(
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(r.Header),
		)
		span := tracer.StartSpan("ping", ext.RPCServerOption(wireContext))
		defer span.Finish()
		r = r.WithContext(opentracing.ContextWithSpan(r.Context(), span))

		w.Write([]byte("pong"))
	})
	http.ListenAndServe(":8080", nil)
}

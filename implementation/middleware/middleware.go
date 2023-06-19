package auapmmiddleware

import (
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"
	"net/http"
)

// AddTraceHeadersToResponse adds the APM trace headers to the response.
//
// Does nothing if there is no transaction in the context.
func AddTraceHeadersToResponse(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if tx := apm.TransactionFromContext(ctx); tx != nil {
			var traceContext = tx.TraceContext()
			headerValue := apmhttp.FormatTraceparentHeader(traceContext)
			w.Header().Set(apmhttp.W3CTraceparentHeader, headerValue)
			if tracestate := traceContext.State.String(); tracestate != "" {
				w.Header().Set(apmhttp.TracestateHeader, tracestate)
			}
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// SetupDiscardTracer globally sets up a discard tracer that does not send any traces.
//
// We override the default tracer because this will affect all APM traces, even those created by
// libraries we include, which often just use the default tracer.
func SetupDiscardTracer() error {
	// if apm is not configured, we use a discardTracer that does not send any traces
	discardTracer, err := apm.NewTracerOptions(apm.TracerOptions{Transport: transport.Discard})
	if err == nil {
		// Set defaultTracer as is also used when starting independent transactions (see scheduler)
		apm.SetDefaultTracer(discardTracer)
	}

	// if there was an error creating the discardTracer we stick with the defaultTracer as a crude backup.
	// The default tracer sends its traces to localhost if it is not configured.

	return err
}

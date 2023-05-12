package auapmclient

import (
	"context"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"net/http"
)

// AddTraceHeadersRequestManipulator is a header manipulator that adds the apm specific trace headers.
//
// Does nothing if there is no transaction in the context.
// For information about the headers and distributed tracing, see https://www.elastic.co/guide/en/apm/guide/8.7/distributed-tracing.html
func AddTraceHeadersRequestManipulator(ctx context.Context, r *http.Request) {
	if tx := apm.TransactionFromContext(ctx); tx != nil {
		traceContext := tx.TraceContext()
		if traceContext.Options.Recorded() {
			if span := apm.SpanFromContext(ctx); span != nil {
				if !span.Dropped() {
					traceContext = span.TraceContext()
				}
			}
		}
		apmhttp.SetHeaders(r, traceContext, false)
	}
}

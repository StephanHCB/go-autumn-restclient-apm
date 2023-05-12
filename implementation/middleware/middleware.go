package auapmmiddleware

import (
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
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

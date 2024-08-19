package auapmclient

import (
	"context"
	"fmt"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"net/http"
)

type apmAutumnContextKey string

var contextKey apmAutumnContextKey = "apmautumn_context"

type RequestTracingImpl struct {
	Wrapped aurestclientapi.Client
}

func New(wrapped aurestclientapi.Client) aurestclientapi.Client {
	return &RequestTracingImpl{
		Wrapped: wrapped,
	}
}

// Perform creates a new APM exit span before calling the wrapped perform and closes it after the wrapped method returns.
//
// Does nothing if there is no transaction in the context.
// For a definition of spans, see https://www.elastic.co/guide/en/apm/guide/current/data-model-spans.html
func (c *RequestTracingImpl) Perform(ctx context.Context, method string, requestUrl string, requestBody interface{}, response *aurestclientapi.ParsedResponse) error {
	ctx, span := prepareApmSpan(ctx, method, requestUrl)
	if span != nil {
		defer span.End()
	}

	err := c.Wrapped.Perform(ctx, method, requestUrl, requestBody, response)

	finalizeApmSpan(span, response.Status, err)
	return err
}

func prepareApmSpan(ctx context.Context, method string, requestUrl string) (context.Context, *apm.Span) {
	isolatedApmContext := context.WithValue(ctx, contextKey, fmt.Sprintf("%s-%s", method, requestUrl))
	tx := apm.TransactionFromContext(isolatedApmContext)

	var span *apm.Span
	if tx != nil {
		traceContext := tx.TraceContext()
		if traceContext.Options.Recorded() {
			// for security reasons we drop unnecessary information by making an incomplete copy
			bodylessReqCopy, err := http.NewRequestWithContext(isolatedApmContext, method, requestUrl, nil)
			if err == nil {
				name := apmhttp.ClientRequestName(bodylessReqCopy)
				span = tx.StartExitSpan(name, "external.http", apm.SpanFromContext(isolatedApmContext))
				if span != nil {
					if !span.Dropped() {
						//put new span into the context to use it for the trace headers instead of the parent trace
						ctx = apm.ContextWithSpan(ctx, span)
						req := apmhttp.RequestWithContext(ctx, bodylessReqCopy)
						span.Context.SetHTTPRequest(req)
					} else {
						span.End()
						span = nil
					}
				}
			}
		}
	}
	return ctx, span
}

func finalizeApmSpan(span *apm.Span, responseStatus int, err error) {
	if err != nil {
		if span != nil {
			span.Context.SetHTTPStatusCode(responseStatus)
		}
	}
}

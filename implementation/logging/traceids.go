package auapmlogging

import (
	"context"
	"go.elastic.co/apm/v2"
)

const TraceIdLogFieldName = "trace.id"
const TransactionIdLogFieldName = "transaction.id"
const SpanIdLogFieldName = "span.id"

// ExtractTraceId returns the trace.id from the TraceContext.
//
// Does nothing if there is no transaction in the context.
func ExtractTraceId(ctx context.Context) string {
	if tx := apm.TransactionFromContext(ctx); tx != nil {
		return tx.TraceContext().Trace.String()
	}
	return ""
}

// ExtractTransactionId returns the transaction.id from the TraceContext.
//
// Does nothing if there is no transaction in the context.
func ExtractTransactionId(ctx context.Context) string {
	if tx := apm.TransactionFromContext(ctx); tx != nil {
		return tx.TraceContext().Span.String()
	}
	return ""
}

// ExtractSpanId returns the span.id from the context.
//
// Does nothing if there is no span in the context.
func ExtractSpanId(ctx context.Context) string {
	if span := apm.SpanFromContext(ctx); span != nil {
		return span.TraceContext().Span.String()
	}
	return ""
}

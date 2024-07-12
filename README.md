# go-autumn-restclient-apm

Adds Elastic APM tracing capabilities based on [elastic/apm-agent-go](https://github.com/elastic/apm-agent-go). Requires
an activated [go.elastic.co/apm/module/apmchiv5/v2](https://pkg.go.dev/go.elastic.co/apm/module/apmchiv5/v2) middleware
to work properly.

## About go-autumn-restclient

It's a rest client that also supports x-www-form-urlencoded.

## About go-autumn-restclient-apm

This library includes the following

- a client wrapper for recording downstream calls in exit spans
- a request manipulator for adding APM trace headers to downstream requests
- a middleware to add APM trace headers to responses
- utilities to extract tracing ids for logging purposes

For the library to work you need to have the
[go.elastic.co/apm/module/apmchiv5/v2](https://pkg.go.dev/go.elastic.co/apm/module/apmchiv5/v2) middleware added to
your router.

## Usage

### Client Wrapper

Change the set-up of your rest client like this:

```
// [...]
apmClient := auapmclient.New(httpClient, circuitBreakerName, maxNumRequestsInHalfOpenState, counterClearingIntervalWhileClosed, timeUntilHalfopenAfterOpen, requestTimeout)
```

You should usually insert the apmClient below the retryer.

### Downstream Trace Headers Request Manipulator

Use the request manipulator function as a parameter while constructing your auresthttpclient.

```
// [...]
client, err := auresthttpclient.New(0, nil, auapmclient.AddTraceHeadersRequestManipulator)
```

### Response Trace Headers Middleware

Just add the middleware to your Router as usual:

```
// [...]
s.Router.Use(auapmmiddleware.AddTraceHeadersToResponse)
```

### APM Logging fields

Use the following extraction methods and the respective field names to add the APM specific fields to your
logs.

| Tracing Id     | Extraction Method                        | Field name constant                      |
|----------------|------------------------------------------|------------------------------------------|
| Transaction Id | `auapmlogging.ExtractTransactionId(ctx)` | `auapmlogging.TransactionIdLogFieldName` |
| Trace Id       | `auapmlogging.ExtractTraceId(ctx)`       | `auapmlogging.TraceIdLogFieldName`       |                                       |
| Span Id        | `auapmlogging.ExtractSpanId(ctx)`        | `auapmlogging.SpanIdLogFieldName`        |

The following code is an example usage with the
[go-autumn-logging-zerolog](https://github.com/StephanHCB/go-autumn-logging-zerolog) middleware.

```
loggermiddleware.AddCustomJsonLogField(auapmlogging.TransactionIdLogFieldName, func(r *http.Request) string {
    return auapmlogging.ExtractTransactionId(r.Context())
})
loggermiddleware.AddCustomJsonLogField(auapmlogging.TraceIdLogFieldName, func(r *http.Request) string {
    return auapmlogging.ExtractTraceId(r.Context())
})
loggermiddleware.AddCustomJsonLogField(auapmlogging.SpanIdLogFieldName, func(r *http.Request) string {
    return auapmlogging.ExtractSpanId(r.Context())
})
```

## Concepts

### Elastic APM

This library heavily depends on the concepts imposed by Elastic APM. For a good start we recommend to read through
the [data model](https://www.elastic.co/guide/en/apm/guide/current/data-model.html) description as well as the
explanation of [distributed tracing](https://www.elastic.co/guide/en/apm/guide/current/apm-distributed-tracing.html) for
the current version of the official documentation.

### Example integration into go-autumn-chi service

The following diagram depicts how this library integrates into the request processing of an existing go-autumn-chi
service.
![visualization_request_processing.svg](docs%2Fvisualization_request_processing.svg)

For scheduled or other not-request-based tasks the following example shows how this library may be utilized to add
Elastic APM tracing. 
![visualization_scheduled_task.svg](docs%2Fvisualization_scheduled_task.svg)


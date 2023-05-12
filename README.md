# go-autumn-restclient-apm

Adds Elastic APM tracing capabilities based on [elastic/apm-agent-go](https://github.com/elastic/apm-agent-go). Requires
an activated [go.elastic.co/apm/module/apmchiv5/v2](https://pkg.go.dev/go.elastic.co/apm/module/apmchiv5/v2) middleware
to work properly.

## About go-autumn

A collection of libraries
for [enterprise microservices](https://github.com/StephanHCB/go-mailer-service/blob/master/README.md) in golang that

- is heavily inspired by Spring Boot / Spring Cloud
- is very opinionated
- names modules by what they do
- unlike Spring Boot avoids certain types of auto-magical behaviour
- is not a library monolith, that is every part only depends on the api parts of the other components
  at most, and the api parts do not add any dependencies.

Fall is my favourite season, so I'm calling it go-autumn.

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
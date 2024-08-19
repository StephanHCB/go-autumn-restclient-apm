package auapmclient

import "net/http"

type ApmRoundTripper struct {
	wrapped http.RoundTripper
}

func NewApmRoundTripper(wrapped http.RoundTripper) *ApmRoundTripper {
	return &ApmRoundTripper{
		wrapped: wrapped,
	}
}

func (vrt *ApmRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	_, span := prepareApmSpan(req.Context(), req.Method, req.RequestURI)
	defer span.End()

	response, err := vrt.wrapped.RoundTrip(req)
	if nil != response {
		finalizeApmSpan(span, response.StatusCode, err)
	}

	return response, err
}

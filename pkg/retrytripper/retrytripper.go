package retrytripper

import (
	"bytes"
	"io"
	"net/http"
)

type ShouldRetryFunc func(*http.Request, *http.Response, error) bool

type Retrytripper struct {
	shouldRetryFunc  ShouldRetryFunc
	realRoundtripper http.RoundTripper
}

func readBody(r *http.Request) ([]byte, error) {
	if r.Body == nil || r.Body == http.NoBody {
		return nil, nil
	}
	return io.ReadAll(r.Body)
	// Do not close the body yet.
}

func (rt *Retrytripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// FIXME: Maybe a body size limit and don't retry huge requests?
	bodyByte, err := readBody(req)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	if bodyByte != nil {
		req.Body = io.NopCloser(bytes.NewReader(bodyByte))
	}
	if rt.realRoundtripper != nil {
		res, err = rt.realRoundtripper.RoundTrip(req)
	} else {
		res, err = http.DefaultTransport.RoundTrip(req)
	}

	shouldRetry := false
	if rt.shouldRetryFunc != nil {
		shouldRetry = rt.shouldRetryFunc(req, res, err)
	}
	if !shouldRetry {
		return res, err
	}

	// Only one retry and no retry delay.
	if bodyByte != nil {
		req.Body = io.NopCloser(bytes.NewReader(bodyByte))
	}
	if rt.realRoundtripper != nil {
		return rt.realRoundtripper.RoundTrip(req)
	} else {
		return http.DefaultTransport.RoundTrip(req)
	}

}

var _ http.RoundTripper = (*Retrytripper)(nil)

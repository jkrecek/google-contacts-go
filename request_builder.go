package contacts

import (
	"io"
	"net/http"
)

type requestBuilder struct {
	client    *Client
	method    string
	path      string
	body      io.Reader
	requestFn func(*http.Request)
	v         interface{}
}

func newRequest(client *Client, method string, path string) *requestBuilder {
	return &requestBuilder{
		client: client,
		method: method,
		path:   path,
	}
}

func (r *requestBuilder) setBody(body io.Reader) *requestBuilder {
	r.body = body
	return r
}

func (r *requestBuilder) setRequestFn(requestFn func(*http.Request)) *requestBuilder {
	r.requestFn = requestFn
	return r
}

func (r *requestBuilder) setInto(v interface{}) *requestBuilder {
	r.v = v
	return r
}

func (r *requestBuilder) Do() error {
	url := r.client.getUrl(r.path)
	req, err := http.NewRequest(r.method, url, r.body)
	if err != nil {
		return err
	}

	req.Header.Set("Gdata-version", "3.0")

	if r.body != nil {
		req.Header.Set("Content-Type", "application/atom+xml")
	}

	if r.requestFn != nil {
		r.requestFn(req)
	}

	resp, err := r.client.native.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return apiError(url, resp.Body)
	} else {
		if r.v != nil {
			return decodeResult(resp.Body, r.v)
		} else {
			return nil
		}

	}
}

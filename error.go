package contacts

import (
	"encoding/xml"
	"fmt"
	"github.com/spkg/bom"
	"io"
	"strings"
)

type jsonErrorResponse struct {
	XMLName xml.Name     `xml:"http://schemas.google.com/g/2005 errors"`
	Errors  []*jsonError `xml:"error,omitempty"`
}

func (e *jsonErrorResponse) Error() string {
	errs := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		errs[i] = err.Error()
	}

	return strings.Join(errs, "\n")
}

type jsonError struct {
	XMLName        xml.Name `xml:"error"`
	Domain         string   `xml:"domain,omitempty"`
	Code           string   `xml:"code,omitempty"`
	InternalReason string   `xml:"internalReason,omitempty"`
}

func (e *jsonError) Error() string {
	return fmt.Sprintf("Json Error: Code `%s`, Domain `%s`, Internal Reason `%s`.", e.Code, e.Domain, e.InternalReason)
}

type ApiError struct {
	url           string
	responseError *jsonErrorResponse
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Request to URL `%d` had an error: %s.\n", e.url, e.responseError)
}

func apiError(url string, rc io.ReadCloser) error {
	r := bom.NewReader(rc)
	dec := xml.NewDecoder(r)

	wrp := new(jsonErrorResponse)
	err := dec.Decode(&wrp)
	if err != nil {
		return fmt.Errorf("GoogleContactsAPI Error parsing went wrong. Error: %s", err)
	}

	return &ApiError{
		url:           url,
		responseError: wrp,
	}
}

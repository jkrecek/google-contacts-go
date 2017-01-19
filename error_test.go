package contacts

import (
	"encoding/xml"
	. "gopkg.in/check.v1"
	"regexp"
	"testing"
)

type ErrorSuite struct{}

var _ = Suite(new(ErrorSuite))

func TestError(t *testing.T) { TestingT(t) }

func (s *ErrorSuite) TestMarshal(c *C) {
	resp := new(jsonErrorResponse)
	resp.Errors = []*jsonError{{
		Domain:         "GData",
		Code:           "invalid",
		InternalReason: "Projection must be one of: full, thin, autocomplete, autocompletepic or property- prefixed: fullsdfa",
	}}

	out, err := xml.Marshal(resp)
	c.Assert(err, IsNil)
	expected := `<errors xmlns="http://schemas.google.com/g/2005">
 <error>
  <domain>GData</domain>
  <code>invalid</code>
  <internalReason>Projection must be one of: full, thin, autocomplete, autocompletepic or property- prefixed: fullsdfa</internalReason>
 </error>
</errors>`

	reSpace := regexp.MustCompile(` +`)
	expected = reSpace.ReplaceAllString(expected, " ")
	re := regexp.MustCompile(`\r?\n ?`)
	expected = re.ReplaceAllString(expected, "")

	c.Assert(string(out), Equals, expected)
}

func (s *ErrorSuite) TestUnmarshal(c *C) {
	raw := `<?xml version="1.0" encoding="UTF-8"?>
<errors xmlns="http://schemas.google.com/g/2005">
 <error>
  <domain>GData</domain>
  <code>invalid</code>
  <internalReason>Projection must be one of: full, thin, autocomplete, autocompletepic or property- prefixed: fullsdfa</internalReason>
 </error>
</errors>`

	into := new(jsonErrorResponse)
	err := xml.Unmarshal([]byte(raw), &into)
	c.Assert(err, IsNil)
	c.Assert(len(into.Errors), Equals, 1)
	c.Assert(into.Errors[0].Code, Equals, "invalid")
	c.Assert(into.Errors[0].Domain, Equals, "GData")
	c.Assert(into.Errors[0].InternalReason, Equals, "Projection must be one of: full, thin, autocomplete, autocompletepic or property- prefixed: fullsdfa")
}

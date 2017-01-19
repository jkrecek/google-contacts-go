package contacts

import (
	"encoding/xml"
	. "gopkg.in/check.v1"
	"regexp"
	"testing"
	"time"
)

const (
	CONTACT_EXAMPLE = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:batch="http://schemas.google.com/gdata/batch" xmlns:gContact="http://schemas.google.com/contact/2008" xmlns:gd="http://schemas.google.com/g/2005" xmlns:openSearch="http://a9.com/-/spec/opensearchrss/1.0/">
 <id>easycore.sync.bridge@gmail.com</id>
 <updated>2017-01-19T08:54:51.668Z</updated>
 <category scheme="http://schemas.google.com/g/2005#kind" term="http://schemas.google.com/contact/2008#contact"/>
 <title type="text">Easy Core's Contacts</title>
 <link rel="alternate" type="text/html" href="https://www.google.com/"/>
 <link rel="http://schemas.google.com/g/2005#feed" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/thin"/>
 <link rel="http://schemas.google.com/g/2005#post" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/thin"/>
 <link rel="http://schemas.google.com/g/2005#batch" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/thin/batch"/>
 <link rel="self" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/thin?max-results=25"/>
 <author>
  <name>Easy Core</name>
  <email>easycore.sync.bridge@gmail.com</email>
 </author>
 <generator version="1.0" uri="http://www.google.com/m8/feeds">Contacts</generator>
 <openSearch:totalResults>2</openSearch:totalResults>
 <openSearch:startIndex>1</openSearch:startIndex>
 <openSearch:itemsPerPage>25</openSearch:itemsPerPage>
 <entry gd:etag="&quot;Q3Y6cDVSLyt7I2A9XR5XGEwNRAA.&quot;">
  <id>http://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/base/56d660e0c8cb34c</id>
  <updated>2017-01-19T09:15:42.818Z</updated>
  <app:edited xmlns:app="http://www.w3.org/2007/app">2017-01-19T09:15:42.818Z</app:edited>
  <category scheme="http://schemas.google.com/g/2005#kind" term="http://schemas.google.com/contact/2008#contact"/>
  <title>Jane Doe</title>
  <link rel="http://schemas.google.com/contacts/2008/rel#photo" type="image/*" href="https://www.google.com/m8/feeds/photos/media/easycore.sync.bridge%40gmail.com/56d660e0c8cb34c"/>
  <link rel="self" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full/56d660e0c8cb34c"/>
  <link rel="edit" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full/56d660e0c8cb34c"/>
  <gd:name>
   <gd:fullName>Jane Doe</gd:fullName>
   <gd:givenName>Jane</gd:givenName>
   <gd:familyName>Doe</gd:familyName>
  </gd:name>
  <gd:email rel="http://schemas.google.com/g/2005#home" address="jane.doe@example.com" primary="true"/>
  <gd:email rel="http://schemas.google.com/g/2005#other" address="jane.doah@example.com"/>
  <gContact:groupMembershipInfo deleted="false" href="http://www.google.com/m8/feeds/groups/easycore.sync.bridge%40gmail.com/base/6"/>
 </entry>

</feed>`
)

type ContactSuite struct{}

var _ = Suite(new(ContactSuite))

func TestContact(t *testing.T) { TestingT(t) }

func (s *ContactSuite) TestMarshal(c *C) {
	contact := new(Contact)
	contact.IdUrl = "http://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/base/5d3bbddc8bae90a3"
	contact.Title = "Frank Doe"
	var err error
	contact.Updated, err = time.Parse(time.RFC3339, "2017-01-18T16:09:09.345Z")
	c.Assert(err, IsNil)
	contact.Emails = []*Email{
		{
			Rel:     "http://schemas.google.com/g/2005#other",
			Address: "frank.doe@example.com",
			Primary: true,
		},
		{
			Rel:     "http://schemas.google.com/g/2005#home",
			Address: "jane.doe@example.com",
			Primary: false,
		},
	}

	out, err := xml.Marshal(contact)
	c.Assert(err, IsNil)
	expected := `<entry xmlns="http://www.w3.org/2005/Atom">
 <id>http://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/base/5d3bbddc8bae90a3</id>
 <updated>2017-01-18T16:09:09.345Z</updated>
 <title>Frank Doe</title>
 <email xmlns="http://schemas.google.com/g/2005" rel="http://schemas.google.com/g/2005#other" address="frank.doe@example.com" primary="true"></email>
 <email xmlns="http://schemas.google.com/g/2005" rel="http://schemas.google.com/g/2005#home" address="jane.doe@example.com"></email>
</entry>
`

	reSpace := regexp.MustCompile(` +`)
	expected = reSpace.ReplaceAllString(expected, " ")
	re := regexp.MustCompile(`\r?\n ?`)
	expected = re.ReplaceAllString(expected, "")

	c.Assert(string(out), Equals, expected)
}

func (s *ContactSuite) TestUnmarshal(c *C) {
	into := new(contactFeed)
	err := xml.Unmarshal([]byte(CONTACT_EXAMPLE), &into)
	c.Assert(err, IsNil)
	c.Assert(len(into.Contacts), Equals, 1)
	c.Assert(into.Contacts[0].Name.FullName, Equals, "Jane Doe")
	c.Assert(len(into.Contacts[0].Emails), Equals, 2)
}

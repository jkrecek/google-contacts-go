package contacts

import (
	"encoding/xml"
	"github.com/davecgh/go-spew/spew"
	. "gopkg.in/check.v1"
	"regexp"
	"testing"
	"time"
)

const (
	CONTACT_EXAMPLE = `<?xml version="1.0" encoding="UTF-8"?>
<feed gd:etag="&quot;SHk7eTVSLyt7I2A9XR5WGUkJRAA.&quot;" xmlns="http://www.w3.org/2005/Atom" xmlns:batch="http://schemas.google.com/gdata/batch" xmlns:gContact="http://schemas.google.com/contact/2008" xmlns:gd="http://schemas.google.com/g/2005" xmlns:openSearch="http://a9.com/-/spec/opensearch/1.1/">
 <id>easycore.sync.bridge@gmail.com</id>
 <updated>2017-02-01T10:02:29.701Z</updated>
 <category scheme="http://schemas.google.com/g/2005#kind" term="http://schemas.google.com/contact/2008#contact"/>
 <title>Easy Core's Contacts</title>
 <link rel="alternate" type="text/html" href="https://www.google.com/"/>
 <link rel="http://schemas.google.com/g/2005#feed" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full"/>
 <link rel="http://schemas.google.com/g/2005#post" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full"/>
 <link rel="http://schemas.google.com/g/2005#batch" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full/batch"/>
 <link rel="self" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full?max-results=25"/>
 <author>
  <name>Easy Core</name>
  <email>easycore.sync.bridge@gmail.com</email>
 </author>
 <generator version="1.0" uri="http://www.google.com/m8/feeds">Contacts</generator>
 <openSearch:totalResults>14</openSearch:totalResults>
 <openSearch:startIndex>1</openSearch:startIndex>
 <openSearch:itemsPerPage>25</openSearch:itemsPerPage>
 <entry gd:etag="&quot;QXo-cTVSLit7I2A9XR5WGUkJRQE.&quot;">
  <id>http://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/base/67e2e50a6ea2db</id>
  <updated>2017-02-01T10:00:50.459Z</updated>
  <app:edited xmlns:app="http://www.w3.org/2007/app">2017-02-01T10:00:50.459Z</app:edited>
  <category scheme="http://schemas.google.com/g/2005#kind" term="http://schemas.google.com/contact/2008#contact"/>
  <title>John Doah</title>
  <link rel="http://schemas.google.com/contacts/2008/rel#photo" type="image/*" href="https://www.google.com/m8/feeds/photos/media/easycore.sync.bridge%40gmail.com/67e2e50a6ea2db"/>
  <link rel="self" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full/67e2e50a6ea2db"/>
  <link rel="edit" type="application/atom+xml" href="https://www.google.com/m8/feeds/contacts/easycore.sync.bridge%40gmail.com/full/67e2e50a6ea2db"/>
  <gd:name>
   <gd:fullName>John Doah</gd:fullName>
   <gd:givenName>John</gd:givenName>
   <gd:familyName>Doah</gd:familyName>
  </gd:name>
  <gd:organization rel="http://schemas.google.com/g/2005#other">
   <gd:orgName>Doah Enterprises</gd:orgName>
   <gd:orgTitle>CEO</gd:orgTitle>
  </gd:organization>
  <gd:email rel="http://schemas.google.com/g/2005#home" address="john.doe@example.com" primary="true"/>
  <gd:email rel="http://schemas.google.com/g/2005#work" address="john.work@example.com"/>
  <gd:phoneNumber rel="http://schemas.google.com/g/2005#mobile">111 222 333</gd:phoneNumber>
  <gd:phoneNumber rel="http://schemas.google.com/g/2005#work">111 333 444</gd:phoneNumber>
  <gd:phoneNumber rel="http://schemas.google.com/g/2005#home">111 444 555</gd:phoneNumber>
  <gd:phoneNumber rel="http://schemas.google.com/g/2005#main">111 555 666</gd:phoneNumber>
  <gd:structuredPostalAddress rel="http://schemas.google.com/g/2005#home">
   <gd:formattedAddress>Street 10
POBox
Vyton, Praha, Praha 11000
Czech republic</gd:formattedAddress>
   <gd:street>Street 10</gd:street>
   <gd:pobox>POBox</gd:pobox>
   <gd:neighborhood>Vyton</gd:neighborhood>
   <gd:city>Praha</gd:city>
   <gd:region>Praha country</gd:region>
   <gd:postcode>11000</gd:postcode>
   <gd:country>Czech republic</gd:country>
  </gd:structuredPostalAddress>
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
	mail1 := NewEmail("frank.doe@example.com", "other", true)
	mail2 := NewEmail("jane.doe@example.com", "home", false)
	contact.Emails = []*Email{&mail1, &mail2}

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
	c.Assert(into.Contacts[0].Name.FullName, Equals, "John Doah")
	c.Assert(len(into.Contacts[0].Emails), Equals, 2)
	c.Assert(into.Contacts[0].Emails[0].Address, Equals, "john.doe@example.com")
	c.Assert(into.Contacts[0].Emails[0].GetType(), Equals, "home")
	c.Assert(into.Contacts[0].Emails[0].Primary, Equals, true)
	c.Assert(len(into.Contacts[0].Addresses), Equals, 1)
	c.Assert(into.Contacts[0].Addresses[0].Street, Equals, "Street 10")
	c.Assert(into.Contacts[0].Addresses[0].POBox, Equals, "POBox")
	c.Assert(into.Contacts[0].Addresses[0].City, Equals, "Praha")
	c.Assert(into.Contacts[0].Addresses[0].Postcode, Equals, "11000")
	c.Assert(into.Contacts[0].Addresses[0].Country, Equals, "Czech republic")

	c.Assert(len(into.Contacts[0].Phones), Equals, 4)
	c.Assert(into.Contacts[0].Phones[0].Number, Equals, "111 222 333")
	c.Assert(into.Contacts[0].Phones[0].GetType(), Equals, "mobile")
	c.Assert(into.Contacts[0].Phones[1].Number, Equals, "111 333 444")
	c.Assert(into.Contacts[0].Phones[1].GetType(), Equals, "work")

	spew.Dump(into)
}

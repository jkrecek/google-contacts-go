package contacts

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	urlIdRegex = regexp.MustCompile(`^.+/([0-9a-f]+)$`)
)

type contactFeed struct {
	XMLName          xml.Name               `xml:"http://www.w3.org/2005/Atom feed"`
	Id               string                 `xml:"id,omitempty"`
	TotalResults     int                    `xml:"http://a9.com/-/spec/opensearchrss/1.0/ totalResults"`
	StartIndex       int                    `xml:"http://a9.com/-/spec/opensearchrss/1.0/ startIndex"`
	ItemsPerPage     int                    `xml:"http://a9.com/-/spec/opensearchrss/1.0/ itemsPerPage"`
	Contacts         []*Contact             `xml:"entry,omitempty"`
	GroupMemberships []*GroupMembershipInfo `xml:"groupMembershipInfo,omitempty"`
}

type Contact struct {
	XMLName xml.Name  `xml:"http://www.w3.org/2005/Atom entry"`
	IdUrl   string    `xml:"id,omitempty"`
	Updated time.Time `xml:"updated,omitempty"`
	Title   string    `xml:"title,omitempty"`
	Name    Name      `xml:"name,omitempty"`
	Emails  []*Email  `xml:"email,omitempty"`
}

func (c *Contact) GetId() (string, error) {
	urlParts := urlIdRegex.FindStringSubmatch(c.IdUrl)
	if len(urlParts) != 2 {
		return "", errors.New("Contact GetId() uses invalid regex")
	}

	return urlParts[1], nil
}

func (c *Contact) SetId(id string) {
	c.IdUrl = fmt.Sprintf("http://www.google.com/m8/feeds/contacts/default/base/%s", id)
}

type Name struct {
	XMLName    xml.Name `xml:"http://schemas.google.com/g/2005 name"`
	FullName   string   `xml:"fullName,omitempty"`
	GivenName  string   `xml:"givenName,omitempty"`
	FamilyName string   `xml:"familyName,omitempty"`
}

func NewName(givenName, familyName string) Name {
	return Name{
		GivenName:  givenName,
		FamilyName: familyName,
		FullName:   fmt.Sprintf("%s %s", givenName, familyName),
	}
}

type Email struct {
	XMLName xml.Name `xml:"http://schemas.google.com/g/2005 email"`
	Rel     string   `xml:"rel,attr"`
	Address string   `xml:"address,attr"`
	Primary bool     `xml:"primary,attr,omitempty"`
}

type GroupMembershipInfo struct {
	XMLName xml.Name `xml:"http://schemas.google.com/contact/2008 groupMembershipInfo"`
	Deleted bool     `xml:"deleted,attr"`
	Href    string   `xml:"href,attr"` // NOTE default `My contacts` group seems to be `http://www.google.com/m8/feeds/groups/{userEmail}/base/6`
}

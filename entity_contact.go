package contacts

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	urlIdRegex = regexp.MustCompile(`^.+/([0-9a-f]+)$`)
)

const (
	schemaRef = "http://schemas.google.com/g/2005"
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
	XMLName      xml.Name                   `xml:"http://www.w3.org/2005/Atom entry"`
	IdUrl        string                     `xml:"id,omitempty"`
	Updated      time.Time                  `xml:"updated,omitempty"`
	Title        string                     `xml:"title,omitempty"`
	Name         Name                       `xml:"name,omitempty"`
	Organization *Organization              `xml:"organization,omitempty"`
	Emails       []*Email                   `xml:"email,omitempty"`
	Phones       []*Phone                   `xml:"phoneNumber,omitempty"`
	Addresses    []*StructuredPostalAddress `xml:"structuredPostalAddress,omitempty"`
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

type rel struct {
	Rel string `xml:"rel,attr"`
}

func (r *rel) GetType() string {
	return strings.Replace(r.Rel, schemaRef+"#", "", -1)
}

func (r *rel) SetType(t string) {
	r.Rel = schemaRef + "#" + t
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
	rel
	Address string `xml:"address,attr"`
	Primary bool   `xml:"primary,attr,omitempty"`
}

func NewEmail(address string, t string, primary bool) Email {
	var r rel
	r.SetType(t)

	return Email{
		Address: address,
		rel:     r,
		Primary: primary,
	}
}

type GroupMembershipInfo struct {
	XMLName xml.Name `xml:"http://schemas.google.com/contact/2008 groupMembershipInfo"`
	Deleted bool     `xml:"deleted,attr"`
	Href    string   `xml:"href,attr"` // NOTE default `My contacts` group seems to be `http://www.google.com/m8/feeds/groups/{userEmail}/base/6`
}

type Phone struct {
	XMLName xml.Name `xml:"http://schemas.google.com/g/2005 phoneNumber"`
	rel
	Number string `xml:",innerxml"`
}

func NewPhone(number string, t string) Phone {
	var r rel
	r.SetType(t)

	return Phone{
		Number: number,
		rel:    r,
	}
}

type Organization struct {
	XMLName xml.Name `xml:"http://schemas.google.com/g/2005 organization"`
	rel
	OrganizationName  string `xml:"orgName,omitempty"`
	OrganizationTitle string `xml:"orgTitle,omitempty"`
}

func NewOrganization(name, title string) Organization {
	var r rel
	r.SetType("other")
	return Organization{
		OrganizationName:  name,
		OrganizationTitle: title,
		rel:               r,
	}
}

type StructuredPostalAddress struct {
	XMLName      xml.Name `xml:"http://schemas.google.com/g/2005 structuredPostalAddress"`
	Rel          string   `xml:"rel,attr"`
	Street       string   `xml:"street,omitempty"`
	POBox        string   `xml:"pobox,omitempty"`
	Neighborhood string   `xml:"neighborhood,omitempty"`
	City         string   `xml:"city,omitempty"`
	Region       string   `xml:"region,omitempty"`
	Postcode     string   `xml:"postcode,omitempty"`
	Country      string   `xml:"country,omitempty"`
}

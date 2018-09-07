package contacts

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
)

func (c *Client) GetContacts() (contacts []*Contact, err error) {
	url := c.getUrl("contacts/default/full?max-results=500")
	var pageContacts []*Contact
	for url != "" {
		pageContacts, url, err = c.getContactsPaged(url)
		contacts = append(contacts, pageContacts...)
	}

	return
}

func (c *Client) getContactsPaged(url string) (contacts []*Contact, nextLink string, err error) {
	feed := new(contactFeed)
	err = newUrlRequest(c, "GET", url).setInto(feed).Do()
	if err != nil {
		return
	}

	if feed.getNextLink() != nil {
		nextLink = feed.getNextLink().Href
	}

	contacts = feed.Contacts
	return
}

func (c *Client) GetContactsInGroup(groupId string) (contacts []*Contact, err error) {
	simpleId := path.Base(groupId)
	allContacts, err := c.GetContacts()
	if err != nil {
		return nil, err
	}

	for _, cnt := range allContacts {
		if cnt.IsMemberOf(simpleId) {
			contacts = append(contacts, cnt)
		}
	}

	return
}

func (c *Client) CreateContact(contact *Contact) (*Contact, error) {
	bts, err := xml.Marshal(contact)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(bts)

	respContact := new(Contact)
	err = newRequest(c, "POST", "contacts/default/full").
		setBody(bodyReader).
		setInto(respContact).
		Do()
	if err != nil {
		return nil, err
	}

	return respContact, nil
}

func (c *Client) UpdateContact(address string, contact *Contact) (*Contact, error) {
	bts, err := xml.Marshal(contact)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(bts)

	respContact := new(Contact)
	if address == "" {
		id, err := contact.GetId()
		if err != nil {
			return nil, err
		}

		address = fmt.Sprintf("contacts/default/full/%s", id)
	}

	err = newRequest(c, "PUT", address).
		setBody(bodyReader).
		setInto(respContact).
		setRequestFn(func(request *http.Request) {
			request.Header.Set("If-None-Match", "DUMMY")
		}).
		Do()

	if err != nil {
		return nil, err
	}

	return respContact, nil
}

func (c *Client) DeleteContact(contactId string) error {
	err := newRequest(c, "DELETE", fmt.Sprintf("contacts/default/full/%s", contactId)).
		setRequestFn(func(request *http.Request) {
			request.Header.Set("If-None-Match", "DUMMY2")
		}).
		Do()
	return err
}

func (c *Client) GetGroups() (groups []*Group, err error) {
	url := c.getUrl("groups/default/full")
	var pageGroups []*Group
	for url != "" {
		pageGroups, url, err = c.getGroupsPaged(url)
		groups = append(groups, pageGroups...)
	}

	return
}

func (c *Client) getGroupsPaged(url string) (groups []*Group, nextLink string, err error) {
	feed := new(groupsFeed)
	err = newUrlRequest(c, "GET", url).setInto(feed).Do()
	if err != nil {
		return
	}

	if feed.getNextLink() != nil {
		nextLink = feed.getNextLink().Href
	}

	groups = feed.Groups
	return
}

func (c *Client) CreateGroup(group *Group) (*Group, error) {
	bts, err := xml.Marshal(group)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(bts)

	respGroup := new(Group)
	err = newRequest(c, "POST", "groups/default/full").
		setBody(bodyReader).
		setInto(respGroup).
		Do()
	if err != nil {
		return nil, err
	}

	return respGroup, nil
}

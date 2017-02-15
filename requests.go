package contacts

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
)

func (c *Client) GetContacts() (contacts []*Contact, err error) {
	feed := new(contactFeed)
	err = c.getRequest("contacts/default/full", feed)
	if err != nil {
		return
	}

	contacts = feed.Contacts
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

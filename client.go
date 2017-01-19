package contacts

import (
	"context"
	"encoding/xml"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

const (
	HOST_URL   = "https://www.google.com"
	FEEDS_PATH = "m8/feeds"
)

type Client struct {
	native    *http.Client
	feedsPath string
}

func NewClient(oauthConfig *oauth2.Config, oauthToken *oauth2.Token) *Client {
	ctx := context.Background()

	return &Client{
		native:    oauthConfig.Client(ctx, oauthToken),
		feedsPath: FEEDS_PATH,
	}
}

func NewClientWithNative(native *http.Client) *Client {
	return &Client{
		native:    native,
		feedsPath: FEEDS_PATH,
	}
}

func (c *Client) SetFeedsPath(feedsPath string) {
	c.feedsPath = feedsPath
}

func (c *Client) getUrl(path string) string {
	return fmt.Sprintf("%s/%s/%s", HOST_URL, c.feedsPath, path)
}

func (c *Client) getRequest(path string, v interface{}) error {
	return newRequest(c, "GET", path).setInto(v).Do()
}

func decodeResult(rc io.ReadCloser, v interface{}) error {
	dec := xml.NewDecoder(rc)
	return dec.Decode(v)
}

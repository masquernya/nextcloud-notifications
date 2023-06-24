package cloud

import (
	"encoding/base64"
	logger "github.com/masquernya/nextcloud-notifications/log"
	"github.com/masquernya/nextcloud-notifications/storage"
	"net/http"
)

var log = logger.New("cloud")

type Cloud struct {
	client *http.Client
}

func New() *Cloud {
	c := &Cloud{
		client: &http.Client{},
	}
	return c
}

func (c *Cloud) GetBasicAuthForAdmin() string {
	username := storage.Get().AdminUsername
	password := storage.Get().AdminPassword
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Cloud) GetBasicAuth() string {
	username := storage.Get().LoginUsername
	password := storage.Get().LoginPassword
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type DavResponseError struct {
	Response *http.Response
	Body     []byte
	Err      error
}

func (e *DavResponseError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "Api response error: " + e.Response.Status + " " + string(e.Body)
}

func NewDavError(response *http.Response, body []byte, err error) error {
	return &DavResponseError{
		Response: response,
		Body:     body,
		Err:      err,
	}
}

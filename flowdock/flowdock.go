// Golang Flowdock client (REST API)
package flowdock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const (
	libraryVersion   = "0.0"
	defaultRestURL   = "https://api.flowdock.com/"
	defaultStreamURL = "https://stream.flowdock.com/"
	tokenRestURL     = "https://%s@api.flowdock.com/"
	tokenStreamURL   = "https://%s@stream.flowdock.com/"
	userAgent        = "go-flowdock/" + libraryVersion
	defaultMediaType = "application/json"
)

// A Client manages communication with the Flowdock API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	RestURL *url.URL

	// Streaming URL for API requests.
	StreamURL *url.URL

	// User agent used when communicating with the Flowdock API.
	UserAgent string

	// Services used for talking to different parts of the Flowdock API.
	Flows         *FlowsService
	Messages      *MessagesService
	Users         *UsersService
	Organizations *OrganizationsService
	Inbox         *InboxService
}

func newClient(httpClient *http.Client, baseURL, streamURL *url.URL) *Client {
	c := &Client{
		client:    httpClient,
		RestURL:   baseURL,
		StreamURL: streamURL,
		UserAgent: userAgent,
	}

	c.Flows = &FlowsService{client: c}
	c.Messages = &MessagesService{client: c}
	c.Inbox = &InboxService{client: c}
	c.Users = &UsersService{client: c}
	c.Organizations = &OrganizationsService{client: c}
	return c
}

// NewClient returns a new Flowdock API client. If a nil httpClient is provided,
// http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the goauth2 library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultRestURL)
	streamURL, _ := url.Parse(defaultStreamURL)
	return newClient(httpClient, baseURL, streamURL)
}

// NewClientWithToken returns a new Flowdock API client instantiated with a
// personal token.  Works the same way as NewClient.
func NewClientWithToken(httpClient *http.Client, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(fmt.Sprintf(tokenRestURL, token))
	streamURL, _ := url.Parse(fmt.Sprintf(tokenStreamURL, token))
	return newClient(httpClient, baseURL, streamURL)
}

func (c *Client) baseRequest(method, urlStr string, baseURL url.URL, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := baseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", defaultMediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the RestURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	return c.baseRequest(method, urlStr, *c.RestURL, body)
}

// NewStreamRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the StreamURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewStreamRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	return c.baseRequest(method, urlStr, *c.StreamURL, body)
}

// Do sends an API request and returns the API response. The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

// An ErrorResponse reports the errors caused by an API request.
//
type ErrorResponse struct {
	Response *http.Response // HTTP response
	Data     []byte         // the error details
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Data)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Data = data
	}
	return errorResponse
}

// addOptions adds the parameters in opt as URL query parameters to s. opt must
// be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Golang Flowdock client (REST API)
package flowdock

import (
	"encoding/json"
	"bytes"
	"net/url"
	"net/http"
	"github.com/google/go-querystring/query"
	"reflect"
	"fmt"
	"io/ioutil"
)

const (
        libraryVersion   = "0.0"
		defaultBaseURL   = "https://api.flowdock.com/"
        userAgent        = "go-flowdock/" + libraryVersion
        defaultMediaType = "application/json"
)

// A Client manages communication with the Flowdock API.
type Client struct {
        // HTTP client used to communicate with the API.
        client *http.Client

        // Base URL for API requests.
        BaseURL *url.URL

        // User agent used when communicating with the Flowdock API.
        UserAgent string

        // Services used for talking to different parts of the Flowdock API.
        Flows *FlowsService
}



// NewClient returns a new Flowdock API client. If a nil httpClient is provided,
// http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the goauth2 library).
func NewClient(httpClient *http.Client) *Client {
        if httpClient == nil {
                httpClient = http.DefaultClient
        }
        baseURL, _ := url.Parse(defaultBaseURL)

        c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
        c.Flows = &FlowsService{client: c}
        return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	    fmt.Println("url:", urlStr)
        rel, err := url.Parse(urlStr)
        if err != nil {
                return nil, err
        }

        u := c.BaseURL.ResolveReference(rel)

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

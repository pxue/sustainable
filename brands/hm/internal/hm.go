package internal

// https://app2.hm.com/hmwebservices/service/article/get-article-by-code/hm-canada/Online/0789407001/en.json

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	// User agent used when communicating with the API.
	UserAgent string

	// options
	opts Options

	common service

	Product     *ProductService
	ProductList *ProductListService
	Supplier    *SupplierService
}

// Options can be used to create a customized client
type Options struct {
	Debug       bool
	country     string
	countryCode string
}

type Option func(*Options) error

// Shopper is an Option to set the Shopper ID.
func Debug(b bool) Option {
	return func(o *Options) error {
		o.Debug = b
		return nil
	}
}

func Country(c, code string) Option {
	return func(o *Options) error {
		o.country = c
		o.countryCode = code
		return nil
	}
}

type service struct {
	client *Client
}

const (
	userAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36`
)

// NewClient returns a new API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide a token that will be sent as part of authHeader.
func NewClient(httpClient *http.Client, options ...Option) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{client: httpClient, UserAgent: userAgent}
	for _, opt := range options {
		if err := opt(&c.opts); err != nil {
			return nil, err
		}
	}

	c.common.client = c

	c.Product = (*ProductService)(&c.common)
	c.ProductList = (*ProductListService)(&c.common)
	c.Supplier = (*SupplierService)(&c.common)
	return c, nil
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if c.opts.Debug {
		fmt.Printf("[hm] %s %s\n", method, u.String())
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	defer func() {
		if c.opts.Debug {
			b, _ := httputil.DumpRequest(req, true)
			fmt.Printf("[hm] %s", string(b))
		}
	}()
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
// TODO: Rate limiting
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		//if e, ok := err.(*url.Error); ok {
		//if url, err := url.Parse(e.URL); err == nil {
		//e.URL = sanitizeURL(url).String()
		//return nil, e
		//}
		//}

		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	// check response status code
	//ResponseStatusCode string `json:"responseStatusCode"`
	//StatusCode

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			if c.opts.Debug {
				b, _ := httputil.DumpResponse(resp, true)
				fmt.Printf("[hm]: %s\n", string(b))
			}
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return resp, err
}

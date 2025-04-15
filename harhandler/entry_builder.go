// Package harhandler provides functionality to build HAR entries from bogdanfinn/fhttp requests
// and responses. It ensures correct extraction of request and response data into the HAR format,
// including body content, headers, cookies, and timing information.
package harhandler

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/Mathious6/harkit/converter"
	"github.com/Mathious6/harkit/harfile"
	http "github.com/bogdanfinn/fhttp"
)

// EntryBuilder builds a HAR entry from a bogdanfinn/fhttp request and optionally a response.
type EntryBuilder struct {
	entry *harfile.Entry
}

// NewEntryWithRequest creates a new EntryBuilder from the given HTTP request and optional cookies.
// It clones the request and preserves the request body to ensure that the original request remains
// usable. The body is read once, stored in memory, and restored on both the original and cloned
// request. Additional cookies are added to the cloned request before it is transformed into a
// HAR-compatible structure.
func NewEntryWithRequest(req *http.Request, additionalCookies []*http.Cookie) (*EntryBuilder, error) {
	clonedReq, _ := cloneRequestPreserveBody(req)

	for _, c := range additionalCookies {
		clonedReq.AddCookie(c)
	}

	harReq, err := converter.FromHTTPRequest(clonedReq)
	if err != nil {
		return nil, err
	}

	return &EntryBuilder{
		entry: &harfile.Entry{
			StartedDateTime: time.Now(),
			Time:            -1,
			Request:         harReq,
			Cache:           &harfile.Cache{},
			Timings: &harfile.Timings{
				Send:    -1,
				Wait:    -1,
				Receive: -1,
			},
		},
	}, nil
}

// AddResponse attaches an HTTP response to the HAR entry and records the time elapsed since the
// request was initiated. This sets the response block and populates the receive timing.
func (b *EntryBuilder) AddResponse(resp *http.Response) error {
	harResp, err := converter.FromHTTPResponse(resp)
	if err != nil {
		return err
	}
	b.entry.Response = harResp

	b.entry.Timings.Receive = float64(time.Since(b.entry.StartedDateTime).Milliseconds())
	b.entry.Time = b.entry.Timings.Total()

	return nil
}

// Build finalizes and returns the HAR entry. If resolveIP is true, it attempts to resolve the
// server's IP address using a DNS lookup based on the request URL.
func (b *EntryBuilder) Build(resolveIP bool) *harfile.Entry {
	if resolveIP && b.entry.Request != nil {
		b.entry.ServerIPAddress = resolveServerIPAddress(b.entry.Request.URL)
	}
	return b.entry
}

// resolveServerIPAddress performs a DNS lookup on the given URL and returns the first resolved
// IP address as a string. Returns an empty string on failure. This is a blocking operation.
func resolveServerIPAddress(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	ipAddrs, err := net.DefaultResolver.LookupIPAddr(context.Background(), parsedURL.Hostname())
	if err != nil || len(ipAddrs) == 0 {
		return ""
	}
	return ipAddrs[0].IP.String()
}

// cloneRequestPreserveBody clones an HTTP request and preserves its body by reading it into memory
// and assigning new readers to both the original and cloned request. This allows for safe reuse of
// both requests without consuming the body multiple times.
func cloneRequestPreserveBody(req *http.Request) (*http.Request, error) {
	if req.Body == nil {
		return req.Clone(req.Context()), nil
	}

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(buf))

	clonedReq := req.Clone(req.Context())
	clonedReq.Body = io.NopCloser(bytes.NewReader(buf))

	return clonedReq, nil
}

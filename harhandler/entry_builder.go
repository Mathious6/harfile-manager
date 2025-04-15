// Package harhandler provides functionality to build HAR entries from bogdanfinn/fhttp requests
// and responses.
package harhandler

import (
	"context"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/Mathious6/harkit/converter"
	"github.com/Mathious6/harkit/harfile"
	http "github.com/bogdanfinn/fhttp"
)

// EntryBuilder builds a HAR entry from a bogdanfinn/fhttp request and optionally a response.
type EntryBuilder struct {
	entry *harfile.Entry
}

// NewEntryWithRequest creates a new EntryBuilder with the given request and cookies. It
// immediately attaches the request as HAR request in the underlying entry and merge cookies (if
// any) into the request headers and cookies.
func NewEntryWithRequest(req *http.Request, cookies []*http.Cookie) (*EntryBuilder, error) {
	harReq, err := converter.FromHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	harReq.Headers = mergeCookieHeader(harReq.Headers, cookies)
	harReq.Cookies = mergeCookies(harReq.Cookies, cookies)

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

// AddResponse attaches an HTTP response to the entry and updates timing information.
func (e *EntryBuilder) AddResponse(resp *http.Response) error {
	harResp, err := converter.FromHTTPResponse(resp)
	if err != nil {
		return err
	}
	e.entry.Response = harResp

	e.entry.Timings.Receive = float64(time.Since(e.entry.StartedDateTime).Milliseconds())
	e.entry.Time = e.entry.Timings.Total()

	return nil
}

// Build finalizes the HAR entry. If resolveIP is true, the server IP address will be resolved and
// stored in the entry.
func (e *EntryBuilder) Build(resolveIP bool) *harfile.Entry {
	if resolveIP && e.entry.Request != nil {
		e.entry.ServerIPAddress = getServerIPAddress(e.entry.Request.URL)
	}
	return e.entry
}

// getServerIPAddress resolves the first IP address for the given URL and returns an empty string if
// resolution fails. This is a blocking call and may take time to resolve.
func getServerIPAddress(reqUrl string) string {
	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		return ""
	}

	ips, err := net.DefaultResolver.LookupIPAddr(context.Background(), parsedUrl.Hostname())
	if err != nil || len(ips) == 0 {
		return ""
	}
	return ips[0].IP.String()
}

// mergeCookies appends new cookies to the existing cookie slice.
func mergeCookies(existing []*harfile.Cookie, new []*http.Cookie) []*harfile.Cookie {
	if len(new) == 0 {
		return existing
	}

	merged := make([]*harfile.Cookie, 0, len(existing)+len(new))
	merged = append(merged, existing...)

	for _, c := range new {
		merged = append(merged, &harfile.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			Expires:  c.RawExpires,
			HTTPOnly: c.HttpOnly,
			Secure:   c.Secure,
		})
	}

	return merged

}

// mergeCookieHeader merges new cookies into an existing "Cookie" header, or adds a new "Cookie"
// header if none exists.
func mergeCookieHeader(existing []*harfile.NVPair, new []*http.Cookie) []*harfile.NVPair {
	if len(new) == 0 {
		return existing
	}

	var cookieHeader strings.Builder
	for i, c := range new {
		if i > 0 {
			cookieHeader.WriteString("; ")
		}
		cookieHeader.WriteString(c.Name + "=" + c.Value)
	}

	for _, header := range existing {
		if strings.EqualFold(header.Name, "Cookie") {
			header.Value = header.Value + "; " + cookieHeader.String()
			return existing
		}
	}

	return append(existing, &harfile.NVPair{
		Name:  "Cookie",
		Value: cookieHeader.String(),
	})
}

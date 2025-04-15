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
// immediately attaches the request as a HAR request in the underlying entry, merging cookies
// into both the request headers and cookie list.
func NewEntryWithRequest(req *http.Request, additionalCookies []*http.Cookie) (*EntryBuilder, error) {
	harReq, err := converter.FromHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	harReq.Headers = mergeCookieHeader(harReq.Headers, additionalCookies)
	harReq.Cookies = mergeCookies(harReq.Cookies, additionalCookies)

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

// Build finalizes the HAR entry. If resolveIP is true, the server IP address will be resolved and
// stored in the entry.
func (b *EntryBuilder) Build(resolveIP bool) *harfile.Entry {
	if resolveIP && b.entry.Request != nil {
		b.entry.ServerIPAddress = resolveServerIPAddress(b.entry.Request.URL)
	}
	return b.entry
}

// resolveServerIPAddress resolves the first IP address for the given URL. Returns an empty
// string if resolution fails. This is a blocking call and may take time to resolve.
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

// mergeCookies appends new cookies into an existing cookie slice.
func mergeCookies(existingCookies []*harfile.Cookie, newCookies []*http.Cookie) []*harfile.Cookie {
	if len(newCookies) == 0 {
		return existingCookies
	}

	combined := make([]*harfile.Cookie, 0, len(existingCookies)+len(newCookies))
	combined = append(combined, existingCookies...)

	for _, c := range newCookies {
		combined = append(combined, &harfile.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Domain:   c.Domain,
			Expires:  c.RawExpires,
			HTTPOnly: c.HttpOnly,
			Secure:   c.Secure,
		})
	}
	return combined
}

// mergeCookieHeader merges new cookies into an existing "Cookie" header, or adds a new "Cookie"
// header if one doesn't already exist.
func mergeCookieHeader(existingHeaders []*harfile.NVPair, newCookies []*http.Cookie) []*harfile.NVPair {
	if len(newCookies) == 0 {
		return existingHeaders
	}

	var mergedCookieValue strings.Builder
	for i, c := range newCookies {
		if i > 0 {
			mergedCookieValue.WriteString("; ")
		}
		mergedCookieValue.WriteString(c.Name + "=" + c.Value)
	}

	for _, hdr := range existingHeaders {
		if strings.EqualFold(hdr.Name, "Cookie") {
			hdr.Value += "; " + mergedCookieValue.String()
			return existingHeaders
		}
	}

	return append(existingHeaders, &harfile.NVPair{
		Name:  "Cookie",
		Value: mergedCookieValue.String(),
	})
}

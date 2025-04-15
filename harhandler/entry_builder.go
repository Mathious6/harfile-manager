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

type EntryBuilder struct {
	entry *harfile.Entry
}

func NewEntry() *EntryBuilder {
	return &EntryBuilder{
		entry: &harfile.Entry{
			StartedDateTime: time.Now(),
			Time:            -1,
			Cache:           &harfile.Cache{},
			Timings: &harfile.Timings{
				Send:    -1,
				Wait:    -1,
				Receive: -1,
			},
		},
	}
}

func (e *EntryBuilder) AddRequest(req *http.Request, cookies []*http.Cookie) error {
	harReq, err := converter.FromHTTPRequest(req)
	if err != nil {
		return err
	}

	// If our client has a cookie jar, we need to merge the cookies
	harReq.Headers = mergeCookieHeader(harReq.Headers, cookies)
	harReq.Cookies = mergeCookies(harReq.Cookies, cookies)

	e.entry.Request = harReq
	return nil
}

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

func (e *EntryBuilder) build(resolveIP bool) *harfile.Entry {
	if resolveIP && e.entry.Request != nil {
		e.entry.ServerIPAddress = getServerIPAddress(e.entry.Request.URL)
	}
	return e.entry
}

func getServerIPAddress(reqUrl string) string {
	url, err := url.Parse(reqUrl)
	if err != nil {
		return ""
	}

	// WARN: This is a blocking call and may take time to resolve
	ipAddress, err := net.DefaultResolver.LookupIPAddr(context.Background(), url.Hostname())
	if err != nil || len(ipAddress) == 0 {
		return ""
	}
	return ipAddress[0].IP.String()
}

func mergeCookies(existing []*harfile.Cookie, new []*http.Cookie) []*harfile.Cookie {
	if len(new) == 0 {
		return existing
	}

	merged := make([]*harfile.Cookie, 0, len(existing)+len(new))
	merged = append(merged, existing...)

	for _, cookie := range new {
		merged = append(merged, &harfile.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			Expires:  cookie.RawExpires,
			HTTPOnly: cookie.HttpOnly,
			Secure:   cookie.Secure,
		})
	}

	return merged

}

func mergeCookieHeader(existing []*harfile.NameValuePair, new []*http.Cookie) []*harfile.NameValuePair {
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

	found := false
	for _, header := range existing {
		if strings.EqualFold(header.Name, "Cookie") {
			header.Value = header.Value + "; " + cookieHeader.String()
			found = true
			break
		}
	}

	if !found {
		existing = append(existing, &harfile.NameValuePair{
			Name:  "Cookie",
			Value: cookieHeader.String(),
		})
	}

	return existing
}

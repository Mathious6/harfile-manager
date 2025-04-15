package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mathious6/harkit/harfile"
	http "github.com/bogdanfinn/fhttp"
)

const (
	ContentLengthKey = "Content-Length"
	ContentTypeKey   = "Content-Type"
	CookieKey        = "Cookie"
	SetCookieKey     = "Set-Cookie"
	LocationKey      = "Location"
)

func convertCookies(cookies []*http.Cookie) []*harfile.Cookie {
	harCookies := make([]*harfile.Cookie, len(cookies))

	for i, cookie := range cookies {
		var expires string
		if !cookie.Expires.IsZero() {
			expires = cookie.Expires.Format(time.RFC3339Nano)
		}

		harCookies[i] = &harfile.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			Expires:  expires,
			HTTPOnly: cookie.HttpOnly,
			Secure:   cookie.Secure,
		}
	}

	return harCookies
}

func convertHeaders(header http.Header, contentLength int64) []*harfile.NameValuePair {
	// By default, client adds Content-Length header later on, so we need to add it here
	// We clone the header to avoid modifying the original one to avoid side effects
	clonedHeader := header.Clone()
	if contentLength > 0 && clonedHeader.Get(ContentLengthKey) == "" {
		clonedHeader.Set(ContentLengthKey, fmt.Sprintf("%d", contentLength))
	}

	harHeaders := make([]*harfile.NameValuePair, 0, len(clonedHeader))
	seen := make(map[string]bool)

	// Used to sort headers in HAR file if needed (e.g. https://github.com/bogdanfinn/tls-client)
	order := clonedHeader.Values(http.HeaderOrderKey)
	for _, name := range order {
		canonical := http.CanonicalHeaderKey(name)
		values := clonedHeader.Values(name)

		if len(values) > 0 {
			for _, value := range values {
				harHeaders = append(harHeaders, &harfile.NameValuePair{Name: canonical, Value: value})
			}
			seen[canonical] = true
		}
	}

	for name, values := range clonedHeader {
		if seen[name] || strings.EqualFold(name, http.HeaderOrderKey) {
			continue
		}
		for _, value := range values {
			harHeaders = append(harHeaders, &harfile.NameValuePair{Name: name, Value: value})
		}
	}

	return harHeaders
}

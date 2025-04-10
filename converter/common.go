package converter

import (
	"net/http"
	"strings"
	"time"

	"github.com/Mathious6/harkit/harfile"
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

func convertHeaders(header http.Header) []*harfile.NameValuePair {
	harHeaders := make([]*harfile.NameValuePair, 0, len(header))

	// Used to sort headers in HAR file if needed (e.g. https://github.com/bogdanfinn/tls-client)
	seen := make(map[string]bool)
	for _, name := range header.Values("Header-Order") {
		if values := header.Values(name); len(values) > 0 {
			for _, value := range values {
				harHeaders = append(harHeaders, &harfile.NameValuePair{Name: name, Value: value})
			}
			seen[http.CanonicalHeaderKey(name)] = true
		}
	}

	for name, values := range header {
		if seen[name] || strings.EqualFold(name, "Header-Order") {
			continue
		}
		for _, value := range values {
			harHeaders = append(harHeaders, &harfile.NameValuePair{Name: name, Value: value})
		}
	}

	return harHeaders
}

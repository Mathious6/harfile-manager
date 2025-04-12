package converter

import (
	"bytes"
	"errors"
	"io"

	"github.com/Mathious6/harkit/harfile"
	http "github.com/bogdanfinn/fhttp"
)

func FromHTTPResponse(resp *http.Response) (*harfile.Response, error) {
	if resp == nil {
		return nil, errors.New("response cannot be nil")
	}

	content, err := buildContent(resp)
	if err != nil {
		return nil, err
	}

	return &harfile.Response{
		Status:      int64(resp.StatusCode),
		StatusText:  http.StatusText(resp.StatusCode),
		HTTPVersion: resp.Proto,
		Cookies:     convertCookies(resp.Cookies()),
		Headers:     convertHeaders(resp.Header),
		Content:     content,
		RedirectURL: locateRedirectURL(resp),
		HeadersSize: -1,
		BodySize:    resp.ContentLength,
		Comment:     "Generated from FromHTTPResponse",
	}, nil
}

func locateRedirectURL(resp *http.Response) string {
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if loc, err := resp.Location(); err == nil {
			return loc.String()
		}
	}
	return ""
}

func buildContent(resp *http.Response) (*harfile.Content, error) {
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(buf))

	return &harfile.Content{
		Size:        int64(len(buf)),
		Compression: 0,
		MimeType:    resp.Header.Get(ContentTypeKey),
		Text:        string(buf),
		Encoding:    "",
	}, nil
}

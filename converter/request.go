package converter

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Mathious6/harkit/harfile"
)

func FromHTTPRequest(req *http.Request) (harfile.Request, error) {
	headers, headersSize := convertHeaders(req)
	postData, bodySize, err := extractPostData(req)
	if err != nil {
		return harfile.Request{}, err
	}

	return harfile.Request{
		Method:      req.Method,
		URL:         req.URL.String(),
		HTTPVersion: req.Proto,
		Cookies:     convertCookies(req.Cookies()),
		Headers:     headers,
		QueryString: convertQueryParams(req.URL),
		PostData:    postData,
		HeadersSize: headersSize,
		BodySize:    bodySize,
		Comment:     "Generated from FromHTTPRequest",
	}, nil
}

func convertCookies(cookies []*http.Cookie) []*harfile.Cookie {
	harCookies := make([]*harfile.Cookie, len(cookies))

	for i, cookie := range cookies {
		harCookies[i] = &harfile.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
	}

	return harCookies
}

func convertHeaders(req *http.Request) ([]*harfile.NameValuePair, int64) {
	harHeaders := make([]*harfile.NameValuePair, 0, len(req.Header))

	// Used to sort headers in HAR file if needed (e.g. https://github.com/bogdanfinn/tls-client)
	seen := make(map[string]bool)
	for _, name := range req.Header.Values("Header-Order") {
		if values := req.Header.Values(name); len(values) > 0 {
			for _, value := range values {
				harHeaders = append(harHeaders, &harfile.NameValuePair{Name: name, Value: value})
			}
			seen[http.CanonicalHeaderKey(name)] = true
		}
	}

	for name, values := range req.Header {
		if seen[name] || strings.EqualFold(name, "Header-Order") {
			continue
		}
		for _, value := range values {
			harHeaders = append(harHeaders, &harfile.NameValuePair{Name: name, Value: value})
		}
	}

	headersSize := int64(0)
	for _, h := range harHeaders {
		headersSize += int64(len(h.Name) + len(h.Value))
	}

	return harHeaders, headersSize
}

func convertQueryParams(u *url.URL) []*harfile.NameValuePair {
	var result []*harfile.NameValuePair

	for key, values := range u.Query() {
		for _, value := range values {
			result = append(result, &harfile.NameValuePair{Name: key, Value: value})
		}
	}

	return result
}

func extractPostData(req *http.Request) (*harfile.PostData, int64, error) {
	if req.Body == nil || req.ContentLength == 0 {
		return nil, int64(req.ContentLength), nil
	}

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, 0, err
	}
	defer req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(buf))

	mimeType := req.Header.Get("Content-Type")
	postData := &harfile.PostData{MimeType: mimeType}
	bodySize := int64(len(buf))

	if strings.HasPrefix(mimeType, "application/x-www-form-urlencoded") {
		text := string(buf)
		pairs := strings.SplitSeq(text, "&")

		for pair := range pairs {
			nv := strings.SplitN(pair, "=", 2)
			if len(nv) == 2 {
				name, value := nv[0], nv[1]
				postData.Params = append(postData.Params, &harfile.Param{Name: name, Value: value})
			}
		}

		return postData, bodySize, nil
	}

	if strings.HasPrefix(mimeType, "multipart/form-data") {
		err := req.ParseMultipartForm(32 << 20) // 32 MB limit
		if err != nil {
			return nil, 0, err
		}

		for name, values := range req.MultipartForm.Value {
			for _, value := range values {
				postData.Params = append(postData.Params, &harfile.Param{Name: name, Value: value})
			}
		}

		for name, files := range req.MultipartForm.File {
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					return nil, 0, err
				}
				defer file.Close()

				content, err := io.ReadAll(file)
				if err != nil {
					return nil, 0, err
				}

				postData.Params = append(postData.Params, &harfile.Param{
					Name:        name,
					FileName:    fileHeader.Filename,
					ContentType: fileHeader.Header.Get("Content-Type"),
					Value:       string(content),
				})
			}
		}

		return postData, bodySize, nil
	}

	postData.Text = string(buf)
	return postData, bodySize, nil
}

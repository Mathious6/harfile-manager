package converter_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/Mathious6/harkit/converter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	METHOD   = http.MethodPost
	URL      = "https://example.com/api?foo=bar"
	PROTOCOL = "HTTP/1.1"

	HEADER_ORDER_KEY = "Header-Order"
	CONTENT_TYPE_KEY = "Content-Type"
	COOKIE_KEY       = "Cookie"

	HEADER1_NAME  = "NaMe1"
	HEADER1_VALUE = "value1"
	HEADER2_NAME  = "nAmE2"
	HEADER2_VALUE = "value2"

	COOKIE_NAME  = "name"
	COOKIE_VALUE = "value"

	URL_CONTENT_TYPE = "application/x-www-form-urlencoded"
	BODY_URL         = "foo=bar"

	JSON_CONTENT_TYPE = "application/json"
	BODY_JSON         = `{"foo":"bar"}`

	PART1_NAME         = "name1"
	PART1_VALUE        = "value1"
	PART2_NAME         = "file"
	PART2_VALUE        = "content"
	PART2_FILENAME     = "test.txt"
	PART2_CONTENT_TYPE = "application/octet-stream"
)

func TestConverter_GivenMethod_WhenConvertingHTTPRequest_ThenMethodShouldBeCorrect(t *testing.T) {
	req := createRequest(t, nil, "")

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Equal(t, METHOD, result.Method, "HAR method <> request method")
}

func TestConverter_GivenURL_WhenConvertingHTTPRequest_ThenURLShouldBeCorrect(t *testing.T) {
	req := createRequest(t, nil, "")

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Equal(t, URL, result.URL, "HAR URL <> request URL")
}

func TestConverter_GivenProtocol_WhenConvertingHTTPRequest_ThenProtocolShouldBeCorrect(t *testing.T) {
	req := createRequest(t, nil, "")

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Equal(t, PROTOCOL, result.HTTPVersion, "HAR protocol <> request protocol")
}

func TestConverter_GivenCookies_WhenConvertingHTTPRequest_ThenCookiesShouldBeCorrect(t *testing.T) {
	req := createRequest(t, nil, "")

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Len(t, result.Cookies, 1, "HAR should contain 1 cookie")
	assert.Equal(t, COOKIE_NAME, result.Cookies[0].Name, "HAR cookie name <> request cookie name")
	assert.Equal(t, COOKIE_VALUE, result.Cookies[0].Value, "HAR cookie value <> request cookie value")
}

func TestConverter_GivenHeaders_WhenConvertingHTTPRequest_ThenHeadersShouldBeCorrectAndOrdered(t *testing.T) {
	req := createRequest(t, nil, "")

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Len(t, result.Headers, 3, "HAR should contain 7 headers")
	assert.Equal(t, HEADER2_NAME, result.Headers[0].Name, "HAR header name <> request header name")
	assert.Equal(t, HEADER2_VALUE, result.Headers[0].Value, "HAR header value <> request header value")
	assert.Equal(t, HEADER1_NAME, result.Headers[1].Name, "HAR header name <> request header name")
	assert.Equal(t, HEADER1_VALUE, result.Headers[1].Value, "HAR header value <> request header value")
}

func TestConverter_GivenURLWithQueryString_WhenConvertingHTTPRequest_ThenQueryStringShouldBeCorrect(t *testing.T) {
	req := createRequest(t, nil, "")

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Len(t, result.QueryString, 1, "HAR should contain 1 query string parameters")
	assert.Equal(t, "foo", result.QueryString[0].Name, "HAR query string name <> request query string name")
	assert.Equal(t, "bar", result.QueryString[0].Value, "HAR query string value <> request query string value")
}

func TestConverter_GivenURLEncodedBody_WhenConvertingHTTPRequest_ThenPostDataShouldBeCorrect(t *testing.T) {
	req := createRequest(t, strings.NewReader(BODY_URL), URL_CONTENT_TYPE)

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Len(t, result.PostData.Params, 1, "HAR should contain 1 post data parameters")
	assert.Empty(t, result.PostData.Text, "HAR should have no post data text")
	assert.Equal(t, "foo", result.PostData.Params[0].Name, "HAR post data name <> request post data name")
	assert.Equal(t, "bar", result.PostData.Params[0].Value, "HAR post data value <> request post data value")
	assert.Equal(t, URL_CONTENT_TYPE, result.PostData.MimeType, "HAR post data mime type <> request post data mime type")
}

func TestConverter_GivenJSONBody_WhenConvertingHTTPRequest_ThenPostDataShouldBeCorrect(t *testing.T) {
	req := createRequest(t, strings.NewReader(BODY_JSON), JSON_CONTENT_TYPE)

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Len(t, result.PostData.Params, 0, "HAR should contain 0 post data parameters")
	assert.NotEmpty(t, result.PostData.Text, "HAR should have post data text")
	assert.Equal(t, BODY_JSON, result.PostData.Text, "HAR post data text <> request post data text")
	assert.Equal(t, JSON_CONTENT_TYPE, result.PostData.MimeType, "HAR post data mime type <> request post data mime type")
}

func TestConverter_GivenMultipartBody_WhenConvertingHTTPRequest_ThenPostDataShouldBeCorrect(t *testing.T) {
	body, contentType := createMultipartBody()
	req := createRequest(t, &body, contentType)

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	assert.Len(t, result.PostData.Params, 2, "HAR should contain 2 post data parameters")
	assert.Empty(t, result.PostData.Text, "HAR should have no post data text")
	assert.Equal(t, PART1_NAME, result.PostData.Params[0].Name, "HAR post data name <> request post data name")
	assert.Equal(t, PART1_VALUE, result.PostData.Params[0].Value, "HAR post data value <> request post data value")
	assert.Equal(t, PART2_NAME, result.PostData.Params[1].Name, "HAR post data name <> request post data name")
	assert.Equal(t, PART2_VALUE, result.PostData.Params[1].Value, "HAR post data value <> request post data value")
	assert.Equal(t, PART2_FILENAME, result.PostData.Params[1].FileName, "HAR post data filename <> request post data filename")
	assert.Equal(t, PART2_CONTENT_TYPE, result.PostData.Params[1].ContentType, "HAR post data content type <> request post data content type")
	assert.Equal(t, contentType, result.PostData.MimeType, "HAR post data mime type <> request post data mime type")
}

func TestConverter_GivenHeaders_WhenConvertingHTTPRequest_ThenHeadersSizeShouldBeCorrect(t *testing.T) {
	req := createRequest(t, nil, "")

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	expectedSize := len(HEADER1_NAME) + len(HEADER1_VALUE)
	expectedSize += len(HEADER2_NAME) + len(HEADER2_VALUE)
	expectedSize += len(COOKIE_KEY) + len(COOKIE_NAME+"="+COOKIE_VALUE)
	assert.Equal(t, int64(expectedSize), result.HeadersSize, "HAR header size <> request header size")
}

func TestConverter_GivenBody_WhenConvertingHTTPRequest_ThenBodySizeShouldBeCorrect(t *testing.T) {
	req := createRequest(t, strings.NewReader(BODY_URL), URL_CONTENT_TYPE)

	result, err := converter.FromHTTPRequest(req)
	require.NoError(t, err)

	expectedSize := int64(len(BODY_URL))
	assert.Equal(t, expectedSize, result.BodySize, "HAR body size <> request body size")
}

// createRequest creates and returns a new HTTP request with the specified body and content type.
// It sets various headers including cookies, custom headers, and content type if provided.
// The function also ensures the request protocol is set and validates the request creation.
func createRequest(t *testing.T, body io.Reader, contentType string) *http.Request {
	req, err := http.NewRequest(METHOD, URL, body)
	require.NoError(t, err)

	req.Proto = PROTOCOL

	req.Header.Add(COOKIE_KEY, COOKIE_NAME+"="+COOKIE_VALUE)
	req.Header.Add(HEADER1_NAME, HEADER1_VALUE)
	req.Header.Add(HEADER2_NAME, HEADER2_VALUE)

	req.Header.Add(HEADER_ORDER_KEY, HEADER2_NAME)
	req.Header.Add(HEADER_ORDER_KEY, HEADER1_NAME)

	if contentType != "" {
		req.Header.Add(CONTENT_TYPE_KEY, contentType)
	}

	return req
}

// createMultipartBody constructs a multipart HTTP request body with predefined fields and file content.
// It returns the body as a bytes.Buffer and the corresponding Content-Type header value.
func createMultipartBody() (body bytes.Buffer, contentType string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	_ = writer.WriteField(PART1_NAME, PART1_VALUE)

	fileWriter, _ := writer.CreateFormFile(PART2_NAME, PART2_FILENAME)
	_, _ = fileWriter.Write([]byte(PART2_VALUE))

	writer.Close()

	return buf, writer.FormDataContentType()
}

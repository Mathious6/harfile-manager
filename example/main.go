package main

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/Mathious6/harkit/harhandler"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

const URL = "https://httpbin.org/post?source=harkit"

func main() {
	handler := harhandler.NewHandler()

	client, _ := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		// tls_client.WithCharlesProxy("127.0.0.1", "8888"), // TODO : do not commit this
	)

	// 1. Form URL-encoded
	form := url.Values{}
	form.Set("name", "Pierre")
	form.Set("role", "developer")
	sendRequest(client, handler, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	fmt.Println("✅ Form URL-encoded request sent")

	// 2. JSON
	jsonBody := `{"name":"Pierre","role":"developer"}`
	sendRequest(client, handler, "application/json", strings.NewReader(jsonBody))
	fmt.Println("✅ JSON request sent")

	handler.Save("example.har")
}

func sendRequest(c tls_client.HttpClient, h *harhandler.HARHandler, contentType string, body io.Reader) {
	req, _ := http.NewRequest(http.MethodPost, URL, body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "harkit-example")
	req.AddCookie(&http.Cookie{Name: "example", Value: "cookie"})

	entry := harhandler.NewEntry()
	_ = entry.AddRequest(req)

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_ = entry.AddResponse(resp)

	h.AddEntry(entry)
}

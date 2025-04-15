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

const URL = "https://httpbin.org"

func main() {
	handler := harhandler.NewHandler()

	client, _ := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCharlesProxy("127.0.0.1", "8888"), // TODO : do not commit this
	)

	// 1. Get with query parameters
	sendGetRequest(client, handler, URL+"/get?name=pierre&role=developer")
	fmt.Println("✅ Parameters sent")

	// 2. Set cookies
	sendGetRequest(client, handler, URL+"/cookies/set?name=pierre&role=developer")
	fmt.Println("✅ Cookies set")

	// 3. Form URL-encoded
	form := url.Values{}
	form.Set("name", "Pierre")
	form.Set("role", "developer")
	sendPostRequest(client, handler, URL+"/post", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	fmt.Println("✅ Form URL-encoded request sent")

	// 4. JSON
	jsonBody := `{"name":"Pierre","role":"developer"}`
	sendPostRequest(client, handler, URL+"/post", "application/json", strings.NewReader(jsonBody))
	fmt.Println("✅ JSON request sent")

	handler.Save("example.har")
}

func sendGetRequest(c tls_client.HttpClient, h *harhandler.HARHandler, URL string) {
	req, _ := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "httpbin.org")
	req.Header.Add("User-Agent", "harkit-example")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	req.Header.Add(http.HeaderOrderKey, "accept")
	req.Header.Add(http.HeaderOrderKey, "host")
	req.Header.Add(http.HeaderOrderKey, "user-agent")
	req.Header.Add(http.HeaderOrderKey, "accept-encoding")

	entry := harhandler.NewEntry()

	urlParsed, _ := url.Parse(URL)
	cookies := c.GetCookies(urlParsed)
	_ = entry.AddRequest(req, cookies)

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_ = entry.AddResponse(resp)

	h.AddEntry(entry)
}

func sendPostRequest(c tls_client.HttpClient, h *harhandler.HARHandler, URL string, contentType string, body io.Reader) {
	req, _ := http.NewRequest(http.MethodPost, URL, body)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Host", "httpbin.org")
	req.Header.Add("User-Agent", "harkit-example")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	req.AddCookie(&http.Cookie{Name: "example", Value: "cookie"})

	req.Header.Add(http.HeaderOrderKey, "accept")
	req.Header.Add(http.HeaderOrderKey, "content-length")
	req.Header.Add(http.HeaderOrderKey, "content-type")
	req.Header.Add(http.HeaderOrderKey, "cookie")
	req.Header.Add(http.HeaderOrderKey, "host")
	req.Header.Add(http.HeaderOrderKey, "user-agent")
	req.Header.Add(http.HeaderOrderKey, "accept-encoding")

	entry := harhandler.NewEntry()

	urlParsed, _ := url.Parse(URL)
	cookies := c.GetCookies(urlParsed)
	_ = entry.AddRequest(req, cookies)

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_ = entry.AddResponse(resp)

	h.AddEntry(entry)
}

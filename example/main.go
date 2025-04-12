package main

import (
	"github.com/Mathious6/harkit/harhandler"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func main() {
	handler := harhandler.NewHandler()

	client, _ := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
		// tls_client.WithCharlesProxy("127.0.0.1", "8888"), // TODO : do not commit this
	)

	entry := harhandler.NewEntry()

	req, _ := http.NewRequest(http.MethodPost, "https://httpbin.org/post?source=harkit", nil)
	req.Header.Add("host", "httpbin.org")
	req.Header.Add("accept-encoding", "gzip, deflate, br")
	req.Header.Add("accept", "*/*")
	req.Header.Add("user-agent", "harkit-example")
	req.Header.Add(http.HeaderOrderKey, "host")
	req.Header.Add(http.HeaderOrderKey, "accept-encoding")
	req.Header.Add(http.HeaderOrderKey, "accept")
	req.Header.Add(http.HeaderOrderKey, "user-agent")
	req.Header.Add(http.HeaderOrderKey, "cookie")
	req.AddCookie(&http.Cookie{Name: "example", Value: "cookie"})
	_ = entry.AddRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_ = entry.AddResponse(resp)

	handler.AddEntry(entry)

	handler.Save("example.har")
}

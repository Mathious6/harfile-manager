package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Mathious6/harkit/converter"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func main() {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithCookieJar(jar),
		tls_client.WithCharlesProxy("127.0.0.1", "8888"), // TODO: REMOVE BEFORE PUSHING
	}

	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	req, _ := http.NewRequest(http.MethodGet, "https://tls.peet.ws/api/all", nil)
	req.Header.Add("accept", "*/*")
	req.Header.Add("user-agent", "harkit-example")
	req.Header.Add(http.HeaderOrderKey, "user-agent")
	req.AddCookie(&http.Cookie{Name: "example", Value: "cookie"})

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	har, err := converter.BuildHAR(req, resp)
	if err != nil {
		panic(err)
	}

	jsonBytes, _ := json.MarshalIndent(har, "", "  ")
	_ = os.WriteFile("main.har", jsonBytes, 0644)
	fmt.Println("✅ HAR file saved as main.har")
}

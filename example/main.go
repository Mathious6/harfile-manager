package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Mathious6/harkit/converter"
)

func main() {
	req, _ := http.NewRequest(http.MethodPost, "https://httpbin.org/post?source=harkit", nil)
	req.Header.Add("user-agent", "harkit-example")
	req.Header.Add("accept", "application/json")
	req.Header.Add(converter.HeaderOrderKey, "accept") // TIPS: Header order is important for some TLS clients.
	req.AddCookie(&http.Cookie{Name: "example", Value: "cookie"})

	resp, _ := http.DefaultClient.Do(req)

	har, err := converter.BuildHAR(req, resp)
	if err != nil {
		panic(err)
	}

	jsonBytes, _ := json.MarshalIndent(har, "", "  ")
	_ = os.WriteFile("main.har", jsonBytes, 0644)
	fmt.Println("✅ HAR file saved as main.har")
}

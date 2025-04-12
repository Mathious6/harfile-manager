package converter

import (
	"context"
	"net"
	"time"

	"github.com/Mathious6/harkit/harfile"
	http "github.com/bogdanfinn/fhttp"
)

func BuildHAR(req *http.Request, resp *http.Response) (*harfile.HAR, error) {
	harReq, err := FromHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	harResp, err := FromHTTPResponse(resp)
	if err != nil {
		return nil, err
	}

	return &harfile.HAR{
		Log: &harfile.Log{
			Version: "1.2",
			Creator: &harfile.Creator{
				Name:    "harkit",
				Version: "0.2.0",
			},
			Entries: []*harfile.Entry{
				{
					StartedDateTime: time.Now(),
					Time:            0,
					Request:         harReq,
					Response:        harResp,
					ServerIPAddress: getServerIPAddress(req.Host),
					Cache:           &harfile.Cache{},
					Timings:         &harfile.Timings{},
				},
			},
		},
	}, nil
}

func getServerIPAddress(host string) string {
	ipAddress, err := net.DefaultResolver.LookupIPAddr(context.Background(), host)
	if err != nil || len(ipAddress) == 0 {
		return ""
	}
	return ipAddress[0].IP.String()
}

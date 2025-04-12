package converter

import (
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
					Cache:           &harfile.Cache{},
					Timings:         &harfile.Timings{},
				},
			},
		},
	}, nil
}

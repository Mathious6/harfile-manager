package harhandler

import (
	"context"
	"net"
	"time"

	"github.com/Mathious6/harkit/converter"
	"github.com/Mathious6/harkit/harfile"
	http "github.com/bogdanfinn/fhttp"
)

type EntryBuilder struct {
	entry *harfile.Entry
}

func NewEntry() *EntryBuilder {
	return &EntryBuilder{
		entry: &harfile.Entry{
			StartedDateTime: time.Now(),
			Cache:           &harfile.Cache{},
			Timings:         &harfile.Timings{},
		},
	}
}

func (e *EntryBuilder) AddRequest(req *http.Request) error {
	harReq, err := converter.FromHTTPRequest(req)
	if err != nil {
		return err
	}
	e.entry.Request = harReq
	return nil
}

func (e *EntryBuilder) AddResponse(resp *http.Response) error {
	harResp, err := converter.FromHTTPResponse(resp)
	if err != nil {
		return err
	}
	e.entry.Response = harResp
	return nil
}

func (e *EntryBuilder) Build() *harfile.Entry {
	e.entry.ServerIPAddress = getServerIPAddress(e.entry.Request.URL)

	return e.entry
}

func getServerIPAddress(url string) string {
	host, _, err := net.SplitHostPort(url)

	ipAddress, err := net.DefaultResolver.LookupIPAddr(context.Background(), host)
	if err != nil || len(ipAddress) == 0 {
		return ""
	}
	return ipAddress[0].IP.String()
}

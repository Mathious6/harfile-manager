package harhandler

import (
	"github.com/Mathious6/harkit"
	"github.com/Mathious6/harkit/harfile"
)

type HandlerOption func(*HARHandler)

type HARHandler struct {
	log *harfile.Log

	resolveIPAddress bool
}

func WithServerIPAddress() HandlerOption {
	return func(h *HARHandler) {
		h.resolveIPAddress = true
	}
}

func NewHandler(opts ...HandlerOption) *HARHandler {
	h := &HARHandler{
		log: &harfile.Log{
			Version: "1.2",
			Creator: &harfile.Creator{
				Name:    "harkit",
				Version: harkit.Version,
			},
			Entries: []*harfile.Entry{},
		},
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *HARHandler) AddEntry(builder *EntryBuilder) {
	entry := builder.Build(h.resolveIPAddress)
	h.log.Entries = append(h.log.Entries, entry)
}

func (h *HARHandler) Save(filename string) error {
	har := &harfile.HAR{Log: h.log}
	return har.Save(filename)
}

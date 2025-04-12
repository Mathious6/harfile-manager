package harhandler

import "github.com/Mathious6/harkit/harfile"

type HARHandler struct {
	log *harfile.Log
}

func NewHandler() *HARHandler {
	return &HARHandler{
		log: &harfile.Log{
			Version: "1.2",
			Creator: &harfile.Creator{
				Name:    "harkit",
				Version: "0.2.0",
			},
			Entries: []*harfile.Entry{},
		},
	}
}

func (h *HARHandler) AddEntry(entry *harfile.Entry) {
	h.log.Entries = append(h.log.Entries, entry)
}

func (h *HARHandler) Save(filename string) error {
	har := &harfile.HAR{Log: h.log}
	return har.Save(filename)
}

package model

type CollectorT struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Local       bool   `json:"local"`
	EbpfFlag    bool   `json:"ebpf"`
	Description string `json:"description"`
	CreatedOn   int64  `json:"createdOn"`
	Token       string `json:"token"`
	//Tenant      string `json:"tenant,omitempty"`
	Account string `json:"account,omitempty"`
}

type CollectorForWeb struct {
	CollectorT
	LastPoll int64 `json:"lastPoll"`
}

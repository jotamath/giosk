package models

type ScanResults struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Service   string `json:"service,omitempty"`
	Banner    string `json:"banner,omitempty"`
	Status    string `json:"status"`
	TimeStamp int64  `json:"timestamp"`
}

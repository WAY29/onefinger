package structs

import "net/http"

// 参考byro07/fwhatweb
type FofaFingerPrint struct {
	Rule_id         string   `json:"rule_id"`
	Level           string   `json:"level"`
	Softhard        string   `json:"softhard"`
	Product         string   `json:"product"`
	Company         string   `json:"company"`
	Category        string   `json:"Category"`
	Parent_category string   `json:"parent_category"`
	Rules           [][]Rule `json:"rules"`
}
type Rule struct {
	Match   string
	Content string
}

type RequestOptions struct {
	Timeout int
}

type FetchResult struct {
	Url          string
	Content      []byte
	Headers      http.Header
	HeaderString string
	Certs        []byte
}

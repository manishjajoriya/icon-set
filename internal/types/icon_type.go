package types

type SearchResponse struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

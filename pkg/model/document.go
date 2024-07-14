package model

type Document struct {
	Url            string `json:"Url,omitempty"`                        
	PubDate        uint64 `json:"PubDate,omitempty"`               
	FetchTime      uint64 `json:"FetchTime,omitempty"`          
	Text           string `json:"Text,omitempty"`                      
	FirstFetchTime uint64 `json:"FirstFetchTime,omitempty"`
}
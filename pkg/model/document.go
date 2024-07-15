package model

type Document struct {
	Url            string `db:"url"`
	PubDate        uint64 `db:"pub_date" `
	FetchTime      uint64 `db:"fetch_time"`
	Text           string `db:"text"`
	FirstFetchTime uint64 `db:"first_fetch_time"`
}

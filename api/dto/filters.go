package dto

type Filters struct {
	HasWebsite  *bool `json:"has_website"`
	HasTwitter  *bool `json:"has_twitter"`
	HasTelegram *bool `json:"has_telegram"`
}

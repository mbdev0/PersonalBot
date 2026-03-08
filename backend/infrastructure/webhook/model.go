package webhook

type Webhook struct {
	Content     any      `json:"content"`
	Embeds      []Embeds `json:"embeds"`
	Attachments []any    `json:"attachments"`
}
type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}
type Author struct {
	Name string `json:"name"`
}
type Thumbnail struct {
	URL string `json:"url"`
}

type Footer struct {
	Text string `json:"text"`
}

type Embeds struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Color     int       `json:"color"`
	Fields    []Fields  `json:"fields"`
	Author    Author    `json:"author"`
	Thumbnail Thumbnail `json:"thumbnail"`
	Footer    Footer    `json:"footer"`
}

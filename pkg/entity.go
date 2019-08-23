package pkg

type FullURL struct {
	URL string `json:"url,omitempty"`
}

type ShortenedURL struct {
	Shortened string `json:"shortened,omitempty"`
}
